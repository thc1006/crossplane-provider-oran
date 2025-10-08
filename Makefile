# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire Makefile.
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

run: ## 在本地執行控制器.
	go run ./adapter/cmd/main.go

##@ Build

build: ## 編譯控制器二進位檔.
	go build -o bin/manager ./adapter/cmd/main.go

docker-build: ## 建置控制器的 Docker 映像.
	docker build -t ${IMG} .

docker-push: ## 推送 Docker 映像到倉庫.
	docker push ${IMG}

##@ Deployment

deploy: ## 部署控制器到 K8s 叢集.
	kustomize build config/default | kubectl apply -f -

undeploy: ## 從 K8s 叢集移除控制器.
	kustomize build config/default | kubectl delete -f -
