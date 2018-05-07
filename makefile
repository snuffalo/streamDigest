BIN_ROOT = bin
BIN_NAME = streamdigest
CLEANABLE_FILES = bin cmd logs models restapi/operations restapi/doc.go restapi/embedded_spec.go restapi/server.go
KUBERNETES_CONFIG = ./streamdigest_deploy.yaml
MAIN_FILE = cmd/$(SERVER_NAME)-server/main.go
SERVER_NAME = streamdigest
SWAGGER_FILE = ./swagger.yml

default: BIN_SUB = local
default: validate generate localbuild move

env: BIN_SUB = scratch
env: validate generate scratchbuild move docker kubernetes

validate:
	swagger validate $(SWAGGER_FILE)
generate:
	swagger generate server -A $(SERVER_NAME) -f $(SWAGGER_FILE)
scratchbuild:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s' -installsuffix cgo -o $(BIN_NAME) $(MAIN_FILE)
localbuild:
	go build -o $(BIN_NAME) $(MAIN_FILE)
move:
	mkdir -p $(BIN_ROOT) && mkdir -p $(BIN_ROOT)/$(BIN_SUB) && mv -f $(BIN_NAME) $(BIN_ROOT)/$(BIN_SUB)/$(BIN_NAME)
docker:
	$(eval $(minikube docker-env))
	docker build -t streamdigest:experimental . -f Dockerfile.streamdigest
kubernetes:
	$(eval $(minikube docker-env))
	kubectl create -f $(KUBERNETES_CONFIG)
clean:
	$(eval $(minikube docker-env))
	rm -rf $(CLEANABLE_FILES)
cleanall:
	$(eval $(minikube docker-env))
	kubectl delete -f $(KUBERNETES_CONFIG) && docker image rm streamdigest:experimental
cleanall: clean