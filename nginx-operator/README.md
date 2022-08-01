# nginx-operator

## Creating nginx-opertor

- Checkout [operator-sdk](https://github.com/operator-framework/operator-sdk)
- cd into operator-sdk and run `make install`
- binaries are created under `./build`
- initialzing a boilerplate:
```bash
../operator-sdk/build/operator-sdk init --domain qinkeith.com --repo github.com/example/nginx-operator
```
- scaffolding the API:
```bash
../operator-sdk/build/operator-sdk create api --group operator --version v1alpha1 --kind NginxOperator --resource --controller
```

