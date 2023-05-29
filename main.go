package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/gokrazy/gokrazy"
)

const rawContainerStoragePath = "/perm/container-storage"

func main() {
	logger.Print("starting up")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	defaultArgs := []string{
		"-e", "TERM=rxvt-unicode",
		"-e", "LANG=C.UTF-8",
		"--tty",
		"--pull", "missing",
		"--device", "/dev/console",
		"--device", "/dev/dri",
		"--device", "/dev/fb0",
		"--device", "/dev/tty",
		"--device", "/dev/tty1",
		"--device", "/dev/vga_arbiter",
		"--device", "/dev/snd",
		"--cap-add", "SYS_TTY_CONFIG",
	}

	runArgs, gokrazyArgs := mergeArgs(defaultArgs, os.Args)
	containerName := gokrazyArgs["name"]

	errChan := make(chan error)
	go run(cancel, errChan, containerName, runArgs...)

	<-ctx.Done()

	logger.Print("exiting signal received")
	cleanup(containerName)

	if err := <-errChan; err != nil {
		logger.Fatalf("fatal error encountered: %v", err)
	}

	logger.Print("gracefully shutting down")
}

func run(cancel context.CancelFunc, errChan chan error, containerName string, args ...string) {
	// Ensure we have an up-to-date clock, which in turn also means that
	// networking is up. This is relevant because podman takes what’s in
	// /etc/resolv.conf (nothing at boot) and holds on to it, meaning your
	// container will never have working networking if it starts too early.
	gokrazy.WaitForClock()

	containerStoragePath := path.Join(rawContainerStoragePath)

	if err := mountVar(containerStoragePath); err != nil {
		logger.Fatal(err)
	}

	cleanup(containerName)

	logger.Printf("issuing %v", append([]string{"podman", "run"}, args...))

	if _, err := podman(context.TODO(), append([]string{"run"}, args...)...); err != nil {
		logger.Printf("error during podman run: %v", err)
		cancel()
		errChan <- fmt.Errorf("error during podman run: %v", err)
	} else {
		errChan <- nil
	}

	return
}

func cleanup(containerName string) {
	logger.Printf("cleaning up existing '%s' container if any", containerName)

	if _, err := podman(context.TODO(), "container", "exists", containerName); err == nil {
		if buf, err := podman(context.TODO(), "stop", containerName); err != nil {
			if !strings.Contains(buf.String(), "no such container") {
				logger.Printf("failed to gracefully stop container %s: %v", containerName, err)
				logger.Printf("killing container %s", containerName)

				if _, err := podman(context.TODO(), "kill", containerName); err != nil {
					logger.Printf("failed to kill container %s: %v", containerName, err)
				}
			}
		}
	}

	if _, err := podman(context.TODO(), "container", "exists", containerName); err == nil {
		// If container exists does not error it means the container exists, so clean it up.
		if _, err := podman(context.TODO(), "rm", containerName); err != nil {
			logger.Printf("failed to rm container %s: %v", containerName, err)
		}
	}
}

// mountVar bind-mounts /perm/container-storage to /var if needed.
// This could be handled by an fstab(5) feature in gokrazy in the future.
func mountVar(storagePath string) error {
	b, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		mountpoint := parts[4]
		if mountpoint == "/var" {
			return nil
		}
	}
	if _, err := os.Stat(storagePath); !os.IsNotExist(err) {
		if err := syscall.Mount(storagePath, "/var", "", syscall.MS_BIND, ""); err != nil {
			return fmt.Errorf("mounting %s to /var: %v", storagePath, err)
		}
	} else {
		if err := syscall.Mount("tmpfs", "/var", "tmpfs", 0, ""); err != nil {
			return fmt.Errorf("mounting tmpfs to /var: %v", err)
		}
	}

	return nil
}

// expandPath returns env, but with PATH= modified or added
// such that both /user and /usr/local/bin are included, which podman needs.
func expandPath(env []string) []string {
	extra := "/user:/usr/local/bin"
	found := false
	for idx, val := range env {
		parts := strings.Split(val, "=")
		if len(parts) < 2 {
			continue // malformed entry
		}
		key := parts[0]
		if key != "PATH" {
			continue
		}
		val := strings.Join(parts[1:], "=")
		env[idx] = fmt.Sprintf("%s=%s:%s", key, extra, val)
		found = true
	}
	if !found {
		const busyboxDefaultPATH = "/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin"
		env = append(env, fmt.Sprintf("PATH=%s:%s", extra, busyboxDefaultPATH))
	}
	return env
}

func podman(ctx context.Context, args ...string) (*bytes.Buffer, error) {
	// Write standard error both to os.Stderr and to a buffer for
	// later consumption.
	buf := bytes.NewBuffer(nil)
	mw := io.MultiWriter(os.Stderr, buf)

	podman := exec.CommandContext(ctx, "/usr/local/bin/podman", args...)
	podman.Env = expandPath(os.Environ())
	podman.Env = append(podman.Env, "TMPDIR=/tmp")
	podman.Stdin = os.Stdin
	podman.Stdout = os.Stdout
	podman.Stderr = mw
	if err := podman.Run(); err != nil {
		return buf, fmt.Errorf("%v: %v", podman.Args, err)
	}
	return buf, nil
}

func mergeArgs(defaultArgs, passedArgs []string) ([]string, map[string]string) {
	var runArgs []string
	gokrazyArgs := make(map[string]string)

	for i, arg := range passedArgs {
		switch arg {
		case "-n", "--name":
			gokrazyArgs["name"] = passedArgs[i+1]
		}
	}

	runArgs = append(defaultArgs, passedArgs[1:]...)

	return runArgs, gokrazyArgs
}
