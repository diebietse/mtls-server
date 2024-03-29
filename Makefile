.PHONY: build
build:
	go build -mod=vendor ./cmd/mtls-server

.PHONY: gofmt
deploy:
	go fmt ./

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: test
test:
	go test ./server

.PHONY: docker
docker:
	docker build -t mtls-server .

.PHONY: lint
lint:
	docker run --rm -it \
		-w /src -v $(shell pwd):/src \
		golangci/golangci-lint:v1.56 golangci-lint run \
		-v -c .golangci.yml

.PHONY: upgrade-vendor
upgrade-vendor:
	go get -u ./...
