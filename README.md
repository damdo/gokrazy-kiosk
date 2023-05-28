## gokrazy-kiosk

gokrazy kiosk allows launching a GUI within a container on gokrazy
and displaying it on a display connected to the host device by leveraging frame buffer passthrough.

This is designed, but not limited to browser kiosks.

You can tailor it to other purposes by following the provided Dockerfile and replacing
the chromium installation with a different GUI application.

### usage

The following is an example on how to launch a browser kiosk (chromium).
This example only works with `amd64` and `arm64` architectures.

```
gok add github.com/gokrazy/iptables
gok add github.com/gokrazy/nsenter
gok add github.com/greenpau/cni-plugins/cmd/cni-nftables-portmap
gok add github.com/greenpau/cni-plugins/cmd/cni-nftables-firewall
gok add github.com/gokrazy/podman

# and finally
gok add github.com/damdo/gokrazy-kiosk
```

Then add under `config.json` within `"PackageConfig"`:

```
        "github.com/damdo/gokrazy-kiosk": {
            "CommandLineFlags": [
                "quay.io/damdo/gokrazy-kiosk-chromium:20230529135304",
                "/usr/bin/chromium",
                "--no-sandbox",
                "--no-first-run",
                "--password-store=basic",
                "--simulate-outdated-no-au='Tue, 31 Dec 2199 23:59:59 GMT'",
                "--noerrdialogs",
                "--disable-infobars",
                "--kiosk",
                "--app='https://play.grafana.org'"
            ]
        },
```

### what is left to do
- [ ] full screen + resolution config need improvements
- [ ] make keyboard work (hard to reliably detect due to gokrazy not having udev)
- [ ] the mouse works but only if it's the only usb device plugged and only if it's there at boot (hard to reliably detect due to gokrazy not having udev)
