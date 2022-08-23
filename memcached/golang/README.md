# golang
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Create the API

- Initialize the project
  ```bash
  operator-sdk init --domain qinkeith.com --repo github.com/qinkeith/operators/memcached/golang
  ```
  - `--repo=` is required when creating a project outsile of `$GOPATH/src` as scaffolding api needs a valid module path
  - [go.mod](./go.mod): used to work with Go modules
  - [Makefilei](./Makefile): Make targets for building/deploying your controller
  - [PROJECT](./PROJECT): Metadata for scaffolding new components
  - [main.go](./main.go): the entrypoint of your controller
    - imports [controller-runtime library](https://pkg.go.dev/sigs.k8s.io/controller-runtime) and it's logging
    - [Schema](https://book.kubebuilder.io/cronjob-tutorial/gvks.html#err-but-whats-that-scheme-thing): provides mappings between Kinds and their corresponding Go types. 

- Scalfold the API
  ```bash
  operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller
  ```
  - Modify [api/v1alpha1/memcached_types.go](./api/v1alpha1/memcached_types.go) to add `size` and  `nodes` to `Spec` and `Status`
  - Update `zz_generated.deepcopy.go` by running
    ```bash
    make generate
    ```
  - Create the CRD manifests at [config/crd/bases/cache.qinkeith.com_memcacheds.yaml](./config/crd/bases/cache.qinkeith.com_memcacheds.yaml) by running
    ```bash
    make manifests
    ```

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/golang:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/golang:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Resources

- [Go Operator Tutorial: memcached-operator](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)
- [Tutorial: Building CronJob](https://book-v2.book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)
- [Explanation of Memcached operator code](https://developer.ibm.com/learningpaths/kubernetes-operators/develop-deploy-simple-operator/deep-dive-memcached-operator-code/) and it's [GitHub repo](https://github.com/IBM/create-and-deploy-memcached-operator-using-go)
- [Initialize and Create an API](https://kubebyexample.com/learning-paths/operator-framework/operator-sdk-go/initialize-and-create-api)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

