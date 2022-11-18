DOCKERCMD = docker

GOCMD = go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

COMMIT = $(shell git rev-parse HEAD)
VERSION = unknown

.PHONY: build
build:
	@$(GOCMD) build \
		-ldflags '-X "github.com/kerraform/kerranamodb//internal/version.Version=$(VERSION)" -X "github.com/kerraform/kerranamodb//internal/version.Commit=$(COMMIT)"' \
		./main.go

TOOLS=\
	github.com/fullstorydev/grpcurl/cmd/grpcurl@latest \
	google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest \
	github.com/bufbuild/buf/cmd/buf@latest

.PHONY: install-tools
install-tools:
	@for tool in $(TOOLS) ; do \
		go install $$tool; \
	done

.PHONY: run
run:
	@$(GOCMD) run \
		-ldflags '-X "github.com/kerraform/kerranamodb//internal/version.Version=$(VERSION)" -X "github.com/kerraform/kerranamodb//internal/version.Commit=$(COMMIT)"' \
		./main.go

.PHONY: run-jaeger
run-jaeger:
	@$(DOCKERCMD) run -d --name jaeger \
		-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
		-p 5775:5775/udp \
		-p 6831:6831/udp \
		-p 6832:6832/udp \
		-p 5778:5778 \
		-p 16686:16686 \
		-p 14268:14268 \
		-p 9411:9411 \
		jaegertracing/all-in-one:1.6

.PHONY: run-01
run-01:
	@LOCK_NODES=http://localhost:18888,http://localhost:18889,http://localhost:18890 \
		HTTP_PORT=8888 GRPC_PORT=18888 \
		go run ./main.go

.PHONY: run-02
run-02:
	@LOCK_NODES=http://localhost:18888,http://localhost:18889,http://localhost:18890 \
		HTTP_PORT=8889 GRPC_PORT=18889 \
		go run ./main.go

.PHONY: run-03
run-03:
	@LOCK_NODES=http://localhost:18888,http://localhost:18889,http://localhost:18890 \
		HTTP_PORT=8890 GRPC_PORT=18890 \
		go run ./main.go
