# Memcached Ansible operator

## Scalfolding the API

- Initialing the project

  ```sh
  operator-sdk init --plugins=ansible --domain qinkeith.com
  ```

- Creating the API

  ```sh
  operator-sdk create api --group cache --version v1alpha1 --kind Memcached --generate-role
  ```

## Update the controller manager

- Add ansible role for the memcached deployment to [roles/memcached/tasks/main.yml](./roles/memcached/tasks/main.yml). This role ensures
  - Memcached deployment exists
  - Deployment size
  Note: the variable `ansible_operator_meta` is retrived from memcached CR, see below in the [memcached CR manifest](#Retrive-details-of-the-CR).
- Set default values in [roles/memcached/defaults/main.yml](./roles/memcached/defaults/main.yml), in our case: `size: 1`
- Add size in [config/samples/cache_v1alpha1_memcached.yaml](./config/samples/cache_v1alpha1_memcached.yaml), `size: 2`

## Configure the image regisrty

- Make changes to Makefile to match docker hub settings
- build and push to docker hub
  
  ```sh
  make docker-build docker-push
  ```

## Create the memcached CR

```sh
kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml

memcached.cache.qinkeith.com/ansible-controller-manager created
```

This creates the memcached CR. Based on the output, you can retrieve details about the CR.

## Retrive details of the CR

```sh
kubectl get memcached/ansible-controller-manager -o yaml

apiVersion: cache.qinkeith.com/v1alpha1
kind: Memcached
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"cache.qinkeith.com/v1alpha1","kind":"Memcached","metadata":{"annotations":{},"name":"ansible-controller-manager","namespace":"default"},"spec":{"size":3}}
  creationTimestamp: "2022-08-21T19:41:17Z"
  generation: 1
  name: ansible-controller-manager
  namespace: default
  resourceVersion: "37220"
  uid: fc50987a-2567-4804-bfde-207aa47f869b
spec:
  size: 3
status:
  conditions:
  ...
```

This CR ensures that the memcached ansible role will be run and a deployment called `ansible-controller-manager-memcached` is created in the default namespace:

```sh
kubectl get deploy

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
ansible-controller-manager-memcached   3/3     3            3           6m9s
```

## Deploy the operator


```sh
make deploy
```

This creates a new namespace and deeployment, in our case `ansible-system` and 'ansible-controller-manager`.

Makefile uses `kustomize` to apply custom configurations and generate manifests from the config/ directory, which are piped to kubectl.

```sh
kustomize build config/default | kubectl apply -f -
```

Namespace and deployment are specified in [config/default/kustomization.yaml](./config/default/kustomization.yaml) and [config/default/manager_config_patch.yaml](./config/default/manager_config_patch.yaml) 

## Verify the custom resource is created

- check memcached CR: `kubectl get memcached/memcached-sample -o yaml`
- There shouold be a new namespace similar to `ansible-system`: ``kubectl get ns`
- A new custom resource should be created: `kubectl api-resource | grep memcached`
- `kubectl get deploy -n ansible-system`

## Update memcached CR

- change size `kubectl patch memcached ansible-controller-manager -p '{"spec":{"size": 2}}' --type=merge`
- run: `kubectl get deploy` to verifiy the changes:
  ```sh
  NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
  ansible-controller-manager-memcached   2/2     2            2           41m
  ```
