# streamDigest
An API for a digest of a livestream

# Requirements (recommend installing in this order)
- go compiler (https://golang.org/doc/install)
- docker (https://docs.docker.com/install/)
- go swagger (https://github.com/go-swagger/go-swagger)
- kubectl (`brew install kubectl`)
- minikube (https://kubernetes.io/docs/getting-started-guides/minikube/#installation)

# Getting Started
- make sure minikube is running locally and your kubectl is pointed at it.
- run `. ./configure.sh` in the shell you intend to run commands from (if anyone knows a better way to do this, please advise!)

- run a `make env` from the root directory of this project. It will
    - Build a GO executable for streamdigest
    - Build a docker container around it
    - Run the docker in a kubernetes deployment behind a kubernetes service (wherever your kubectl is pointing)
    - Run a swagger ui docker container for documentation

- run a `make` from the root directory to get a plain executable capable of running on your machine (for simple debugging)

# Doin' Stuff
(Note: Nothing is implemented currently in the web app, only the framework)
* if you ran `make env`, you can `echo $(minikube service streamdigest --url)` to get the url to reach the app
* if you ran `make env`, you can `echo $(minikube service streamdigest-docs --url))` to get the url to reach the swagger ui
* if you ran `make`, you can now start the exectuable found in `bin/local` or debug it
