.PHONY: docker-build
IMAGE ?= ghcr.io/your-org/matrix-bumble:latest

docker-build:
	test -n "$$BBCTL_AUTH_JSON" || (echo "BBCTL_AUTH_JSON is required (do not commit it). Example: export BBCTL_AUTH_JSON='{\"device_id\":\"...\",\"environments\":{...}}'"; exit 1)
	docker buildx build --platform linux/amd64,linux/arm64 --build-arg BBCTL_AUTH_JSON="$$BBCTL_AUTH_JSON" -t $(IMAGE) . --push

.PHONY: proxy
proxy:
	bbctl proxy -r registration.yaml

.PHONY: run
run:
	go run . -r registration.yaml
