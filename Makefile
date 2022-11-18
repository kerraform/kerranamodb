DOCKERCMD = docker

IMG ?= karranamodb:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOCMD = go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

DEV_DIR   := $(shell pwd)/dev
BIN_DIR   := $(DEV_DIR)/bin
TOOLS_DIR := $(DEV_DIR)/tools
BUILD_TOOLS := cd $(TOOLS_DIR) && go build -o

# Etcd
ETCD_VERSION           := 3.6.0-alpha.0
ETCD := $(abspath $(BIN_DIR)/etcd)-$(ETCD_VERSION)

# Kind
KIND_VERSION      := 0.14.0
KIND              := $(abspath $(BIN_DIR)/kind)-$(KIND_VERSION)
KIND_CLUSTER_NAME := karranamodb
KIND_MANIFEST     := $(abspath $(DEV_DIR)/kind/cluster.yaml)

# Kubernetes
KUBERNETES_VERSION     := 1.24.1
KUBE_APISERVER_VERSION := 1.24.1
KUBE_APISERVER := $(abspath $(BIN_DIR)/kube-apiserver)-$(KUBE_APISERVER_VERSION)
KUBEBUILDER_VERSION := 3.7.0
KUBEBUILDER         :=  $(abspath $(BIN_DIR)/kubebuilder)

# Scaffold
SKAFFOLD_VERSION := 1.39.1
SKAFFOLD       := $(abspath $(BIN_DIR)/skaffold)-$(SKAFFOLD_VERSION)

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

etcd: $(ETCD)
$(ETCD):
	@curl -sSL "https://github.com/etcd-io/etcd/releases/download/v$(ETCD_VERSION)/etcd-v$(ETCD_VERSION)-$(GOOS)-$(GOARCH).tar.gz" | tar -C /tmp -xzv etcd-v$(ETCD_VERSION)-$(GOOS)-$(GOARCH)/etcd
	@mv /tmp/etcd-v$(ETCD_VERSION)-$(GOOS)-$(GOARCH)/etcd $(ETCD)
	@chmod +x $(ETCD)
	@cp $(ETCD) $(BIN_DIR)/etcd

kind: $(KIND)
$(KIND):
	@curl -Lso $(KIND) https://github.com/kubernetes-sigs/kind/releases/download/v$(KIND_VERSION)/kind-$(GOOS)-$(GOARCH)
	@chmod +x $(KIND)
	@cp $(KIND) $(BIN_DIR)/kind

.PHONY: kind-cluster
kind-cluster: $(KIND) $(KUBECTL)
	@$(KIND) delete cluster --name $(KIND_CLUSTER_NAME)
	@$(KIND) create cluster --name $(KIND_CLUSTER_NAME) --config $(KIND_MANIFEST)
	@make kind-manifest
	@make kind-image

kube-apiserver: $(KUBE_APISERVER)
$(KUBE_APISERVER):
	@curl -sSL "https://dl.k8s.io/v$(KUBE_APISERVER_VERSION)/kubernetes-server-$(GOOS)-$(GOARCH).tar.gz" | tar -C /tmp -xzv kubernetes/server/bin/kube-apiserver
	@mv /tmp/kubernetes/server/bin/kube-apiserver $(KUBE_APISERVER)
	@chmod +x $(KUBE_APISERVER)
	@cp $(KUBE_APISERVER) $(BIN_DIR)/kube-apiserver

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

.PHONY: run-dev
run-dev: $(SKAFFOLD)
	@kubectl config use-context kind-$(KIND_CLUSTER_NAME)
	@PATH=$${PWD}/dev/bin:$${PATH} $(SKAFFOLD) dev --tail --filename=./dev/skaffold/skaffold.yaml

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

.PHONY: run-service-discovery
run-service-discovery:
	@LOCK_SERVICE_DISCOVERY_ENDPOINT=localhost \
		LOCK_SERVICE_DISCOVERY_NODE_COUNT=4 \
		LOCK_SERVICE_DISCOVERY_TIMEOUT=3 \
		HTTP_PORT=8888 GRPC_PORT=18888 \
		go run ./main.go

.PHONY: run-01
run-01:
	@LOCK_NODES=http://localhost:18889,http://localhost:18890,http://localhost:18891 \
		HTTP_PORT=8888 GRPC_PORT=18888 \
		go run ./main.go

.PHONY: run-02
run-02:
	@LOCK_NODES=http://localhost:18888,http://localhost:18890,http://localhost:18891 \
		HTTP_PORT=8889 GRPC_PORT=18889 \
		go run ./main.go

.PHONY: run-03
run-03:
	@LOCK_NODES=http://localhost:18888,http://localhost:18889,http://localhost:18891 \
		HTTP_PORT=8890 GRPC_PORT=18890 \
		go run ./main.go

.PHONY: run-04
run-04:
	@LOCK_NODES=http://localhost:18888,http://localhost:18889,http://localhost:18890 \
		HTTP_PORT=8891 GRPC_PORT=18891 \
		go run ./main.go

skaffold: $(SKAFFOLD)
$(SKAFFOLD):
	@curl -Lso $(SKAFFOLD) https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(GOOS)-$(GOARCH)
	@chmod +x $(SKAFFOLD)
	