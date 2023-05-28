# ==============================================================================
# CLASS NOTES
#
# Kind
#   For full Kind v0.18 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.18.0

run:
	@go run app/services/shortener-api/main.go | go run app/tooling/logfmt/main.go -service=shortener-api

run-help:
	go run app/services/shortener-api/main.go --help

tidy:
	go mod tidy

dev-brew-common:
	brew update
	brew tap hashicorp/tap
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list vault || brew install vault
	brew list staticcheck || brew install staticcheck

dev-brew-arm64: dev-brew-common
	brew list datawire/blackbird/telepresence-arm64 || brew install datawire/blackbird/telepresence-arm64

# ========================================
# Running tests within the local computer

test:
	CGO_ENABLED=0 go test -count=1 ./...
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...
	# govulncheck ./...

# ========================================
# Building containers

# ${shell git rev-parse --short HEAD}
VERSION := 0.1.0

all: shortener

shortener:
	docker build \
		-f zarf/docker/dockerfile.shortener-api \
		-t shortener-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.


# ==============================================================================
# Install dependencies

GOLANG       := golang:1.20
ALPINE       := alpine:3.17
KIND         := kindest/node:v1.26.3
POSTGRES     := postgres:15-alpine
VAULT        := hashicorp/vault:1.13
ZIPKIN       := openzipkin/zipkin:2.24
TELEPRESENCE := docker.io/datawire/tel2:2.13.1

# ========================================

dev-docker:
	# docker pull $(GOLANG)
	# docker pull $(ALPINE)
	# docker pull $(KIND)
	# docker pull $(POSTGRES)
	# docker pull $(VAULT)
	# docker pull $(ZIPKIN)
	docker pull $(TELEPRESENCE)

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest

# ========================================
# Running from k8s/kind

KIND_CLUSTER_NAME := zmoog-starter-cluster

dev-up-local:
	kind create cluster \
	  --image kindest/node:v1.26.3@sha256:61b92f38dff6ccc29969e7aa154d34e38b89443af1a2c14e6cfbd2df6419c66f \
	  --name $(KIND_CLUSTER_NAME) \
	  --config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down-local:
	telepresence quit -s 
	kind delete cluster --name $(KIND_CLUSTER_NAME)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-load:
	# cd zarf/k8s/dev/shortener; kustomize edit set image shortener-api-image=shortener-api:$(VERSION)
	kind load docker-image shortener-api:$(VERSION) --name $(KIND_CLUSTER_NAME)

dev-apply:
	kustomize build zarf/k8s/dev/shortener | kubectl apply -f -
	kubectl wait pods --namespace=shortener-system --selector app=shortener --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment shortener --namespace=shortener-system

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply	

dev-describe:
	kubectl describe nodes
	kubectl describe svc

dev-describe-deployment:
	kubectl describe deployment --namespace=shortener-system shortener

dev-describe-shortener:
	kubectl describe pod --namespace=shortener-system -l app=shortener

# ------------------------------------------------------------------------------

dev-logs:
	kubectl logs --namespace=shortener-system -l app=shortener --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=shortener-api


# ------------------------------------------------------------------------------

dev-tel-setup:
	kind load docker-image $(TELEPRESENCE) --name $(KIND_CLUSTER_NAME)
	telepresence --context=kind-$(KIND_CLUSTER_NAME) helm install

dev-tel-connect: dev-tel-setup
	telepresence --context=kind-$(KIND_CLUSTER_NAME) connect

# ------------------------------------------------------------------------------

metrics-local:
	expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

metrics-view:
	expvarmon -ports="shortener-service.shortener-system.svc.cluster.local:4000" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

metrics-view-sc:
	expvarmon -ports="shortener-service.shortener-system.svc.cluster.local:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

load:
	hey -m GET -c 100 -n 10000 "http://shortener-service.shortener-system.svc.cluster.local:3000/status"

# ------------------------------------------------------------------------------

status:
	curl -il http://shortener-service.shortener-system.svc.cluster.local:3000/status


# RSA Keys
#   To generate a priva/public key pair, run:
#   $ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
#   $ openssl rsa -pubout -in private.pem -out public.pem
jwt:
	go run app/scratch/jwt/main.go
