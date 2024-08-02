all: build-push

IMAGE = quay.io/damdo/gokrazy-kiosk-chromium:$(shell date +%Y%m%d%H%M%S)

build-push-docker:
	echo "docker: creating multi-arch container image: $(IMAGE) ..."
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(IMAGE) .

build-push-podman:
	echo "podman: creating multi-arch container image: $(IMAGE) ..."
	podman manifest create $(IMAGE)
	podman build --platform linux/amd64,linux/arm64 --manifest $(IMAGE)  .
	podman manifest push $(IMAGE)
