BIN_ROOT = bin
BIN_NAME = protocbas
CLEANABLE_FILES = bin cmd logs models restapi/restapi/operations  restapi/doc.go restapi/embedded_spec.go restapi/server.go
MAIN_FILE = cmd/$(SERVER_NAME)-server/main.go
SERVER_NAME = protocbas
SWAGGER_FILE = ./swagger.yml

default: BIN_SUB = local
default: general

test: BIN_SUB = linux_amd64
test: export GOOS=linux
test: export GOARCH=amd64
test: general docker kubernetes

run: test

general: validate generate gobuild move

validate:
	swagger validate $(SWAGGER_FILE)
generate:
	swagger generate server -A $(SERVER_NAME) -f $(SWAGGER_FILE) && mkdir -p $(BIN_ROOT) && mkdir -p $(BIN_ROOT)/$(BIN_SUB)
gobuild:
	go build -o $(BIN_NAME) $(MAIN_FILE)
move:
	mv -f $(BIN_NAME) $(BIN_ROOT)/$(BIN_SUB)/$(BIN_NAME)
docker:
	$(eval $(minikube docker-env))
	docker build -t protocbas:experimental . -f Dockerfile.protocbas
kubernetes:
	kubectl create -f protocbas_deploy.yaml
clean:
	rm -rf $(CLEANABLE_FILES) && docker system prune -f
cleanall:
	clean
	kubernetes destroy -t protocbas_deploy.yaml
