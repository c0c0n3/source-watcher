
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/c0c0n3/osmops:latest
# Produce CRDs that work back to Kubernetes 1.16
CRD_OPTIONS ?= crd:crdVersions=v1

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd
#	kustomize build config/crd | kubectl apply -f -
#   ^ potentially harmful. what if you're connected to the wrong cluster?!

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd
#	kustomize build config/crd | kubectl delete -f -
#   ^ potentially harmful. what if you're connected to the wrong cluster?!

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image source-watcher=${IMG}
	kustomize build config/default
#	kustomize build config/default | kubectl apply -f -
#   ^ potentially harmful. what if you're connected to the wrong cluster?!

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=source-reader webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
# TODO: get rid of this mess. Use GH actions when moving to own repo.
# NOTE. docker login stores the token in osxkeychain, logout removes it from
# there. Not sure why GH recommends using their tool (gh), since as you can
# see the token is stored unencrypted and with no password protection!
docker-push:
	grep oauth_token ~/.config/gh/hosts.yml | sed 's/.*oauth_token: //' | docker login ghcr.io -u c0c0n3 --password-stdin
	docker push ${IMG}
	docker logout ghcr.io

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.5.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
