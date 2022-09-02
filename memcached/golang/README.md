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
  - Generate code for DeepCopy:. [zz_generated.deepcopy.go](api/v1alpha1/zz_generated.deepcopy.go) by running [make generate](./Makefile#L93)
  - Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects for the CRD at [config/crd/bases/cache.qinkeith.com_memcacheds.yaml](./config/crd/bases/cache.qinkeith.com_memcacheds.yaml) by running [make manifests](./Makefile#L89)
    - The controller needs certain RBAC permissions to interact with the resources it manages. These are specified via RBAC markers like the following:
      ```marker
      //+kubebuilder:rbac:groups=cache.qinkeith.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
      //+kubebuilder:rbac:groups=cache.qinkeith.com,resources=memcacheds/status,verbs=get;update;patch
      //+kubebuilder:rbac:groups=cache.qinkeith.com,resources=memcacheds/finalizers,verbs=update
      ```
       in controllers/memcached_controller.go. `make manifest` creates the ClusterRole manifest at [config/rbac/role.yaml](./config/rbac/role.yaml) 
  
  Note, both `make generate` and `make manifests` will call [controller-gen](https://github.com/kubernetes-sigs/controller-tools) utility.

- Implement the Controller
  - Reference the instance you want to observe. In our case, it's the Memcached object definded in [api/v1alpha1/memcached_types.go](https://github.com/qinkeith/operators/blob/main/memcached/golang/api/v1alpha1/memcached_types.go#L43). This can be achieved by importing the Memcached CRD from the `cachev1alpha1` object:
    
    ```golang
      import (
        ...
        cachev1alpha1 "github.com/qinkeith/operators/memcached/golang/api/v1alpha1"
      )	
    ```
   
    `cachev1alpha1.<Object>{}` can be used to reference any of the defined objects within that memcached_types.go. For example:

    ```golang
    memcached := &cachev1alpha1.Memcached{}
    ``` 
  
  - The [Reconcile](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile) method - The Reconcile [method](https://go.dev/tour/methods/1) is 
  a function with MemcachedReconciler as it's receiver:

  ```golang
  type MemcachedReconciler struct {
        client.Client
        Scheme *runtime.Scheme
  }

  func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
  ...
  ```
  
  In Go, a function which takes a receiver is usually called a method of the type (receiver). In our case, Reconcile is a method of MemcachedReconciler. 
  We can make calls such as:

  ```golang
  dep := r.deploymentForMemcached(memcached)
  ```

  Reconcile method is responsible for enforcing the 
  desired CR state on the actual state of the system. It runs each time an event occurs on a watched CR or resource, and will return some value 
  depending on whether those states match or not.

    In this way, every Controller has a Reconciler object with a Reconcile() method that implements the reconcile loop. The reconcile loop is passed 
    the Request argument which is a Namespace/Name key used to lookup the primary resource object, Memcached, from the cache: 

      ```golang
      func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
      ```

      This function expects:
  
    - [Context](https://go.dev/blog/context): The context carries a deadline, a cancellation signal, and other values across API boundaries. The context       takes into account the identity of the end user, auth tokens, and the request's deadline. To view your current context:
        
        ```bash
        kubectl config view
        ```

    - [Request](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile#Request): Request contains the information necessary to reconcile a     Kubernetes object. This includes the information to uniquely identify the object - its Name and Namespace.

    The following are a few possible return options for a Reconciler:
    - with error:
      
      ```golang
      return ctrl.Result{}, err
      ```

    - without error, but reqeue request:
      
      ```golang
      return ctrl.Result{Reqeue: true}, nil 
      ```

    - to stop the Reconcile:
      
      ```golang
      return ctrl.Result{}, nil 
      ```

  - The client [Reader](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.7.0/pkg/client/interfaces.go#L48) interface. Reader knows how to read and list Kubernetes objects. 
    - The [Get](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.Get) function - Use it to confirm that the observed resource, Memcached in our case, is defined in the namespace:

      ```golang
      memcached := &cachev1alpha1.Memcached{}
      err := r.Get(ctx, req.NamespacedName, memcached)
      ```

      We aleady know the first 2 parameters, `context` and `request`. The `req` struct contains the `NamespacedName` which is the name and the namespace 
      of the object to reconcile. The object must be a struct pointer so that memcached can be updated with the response returned by the server. 
      In our case, that is the memcached object which must be a struct pointer so that memcached can be updated with the content returned by the Server.

    - The [List](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.List) function - List function can be used to retrieve a list all    child objects in a given namespace and list options. 
      In our case, we use List to retrieve all [memcached pods](./controllers/memcached_controller.go#L115). On a successful call, `Items` field in the list 
      will be populated with the result returned from the server:
  
      ```golang
      podNames := getPodNames(podList.Items)
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
- [Implement the controller](https://book-v2.book.kubebuilder.io/cronjob-tutorial/controller-implementation.html)

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

