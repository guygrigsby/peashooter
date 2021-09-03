version = 0.0.1
image = peashooter
registry = docker.io/guygrigsby
build = $(image):$(version)
test:
	go test ./... -v

run: build
	@docker run -it $(registry)/$(build)

.PHONY: build
build:
	@echo "Building $(build)..."
	@docker build --rm=true --no-cache=true --pull=true -t $(build) .
	@docker tag $(build) $(registry)/$(build)

.PHONY: release
release: build
	@echo "Releasing $(build)..."
	@docker push $(registry)/$(build)
.PHONY: run test
