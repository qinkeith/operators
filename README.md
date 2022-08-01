# Operators

## Building operator-sdk on Ubuntu

- Checkout [operator-sdk](https://github.com/operator-framework/operator-sdk)
  ```bash
  https://github.com/operator-framework/operator-sdk
  ```
- cd into operator-sdk and run `make install`
- binaries are created under `./build`

## Building operator-sdk on WLS Ubuntu

- Install Ubuntu 22.04 LTS from Microsoft app store
- Run `suso apt upgrade`
- Install go: `sudo apt install golang-go`
- Install make: `sudo apt install make`
- Checkout operator-sdk `git clone git@github.com:operator-framework/operator-sdk.git`

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
