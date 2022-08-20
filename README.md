# Operators

## Setting up operator-sdk on Ubuntu

- Checkout [operator-sdk](https://github.com/operator-framework/operator-sdk)
  ```bash
  git clone https://github.com/operator-framework/operator-sdk
  ```
## Setting up operator-sdk on WLS Ubuntu

- Install Ubuntu 22.04 LTS from Microsoft app store
- Run `sudo apt upgrade`
- Install go: `sudo apt install golang-go`
- Install make: `sudo apt install make`
- Checkout operator-sdk `git clone https://github.com/operator-framework/operator-sdk.git`

## Building operator-sdk

- cd into operator-sdk and run `make`
- binaries are created under `./build`
- run `make install`
- installs to `$GOPATH/bin`

## Setting up minikube for WSL ubuntu
- Go to Docker Desktop->Settings->Resources->WSL Integration
- Enable the desired distro
- Download minikube: `curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64`
- `chmod +x ./minikube`
- `sudo usermod -aG docker $USER && newgrp docker`
- Run: `sudo usermod -aG docker $USER && newgrp docker`
- Start `./minikube start`

## Creating nginx-operator
- mkdir nginx-operator and cd into it
- initialzing a boilerplate using `operator init`:
```bash
../operator-sdk/build/operator-sdk init --domain qinkeith.com --repo github.com/example/nginx-operator
```
- scaffolding the API using `operator api`:
```bash
../operator-sdk/build/operator-sdk create api --group operator --version v1alpha1 --kind NginxOperator --resource --controller
```
