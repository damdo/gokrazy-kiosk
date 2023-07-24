all: build-push

build-push:
	docker buildx build --push --platform linux/amd64,linux/arm64 -t quay.io/damdo/gokrazy-kiosk-chromium:$$(date +%Y%m%d%H%M%S) .
