version = 0.0.1
image = peashooter
registry = docker.io/guygrigsby
build = $(image):$(version)
dev:
	go run cmd/main.go
test:
	go test ./... -v

run: build
	$GOPATH/src/github.com/k4s/phantomgo/phantomgojs
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
