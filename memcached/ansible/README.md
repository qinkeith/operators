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

- Add ansible role for the memcached deployment to `roles/memcached/tasks/main.yml`. This role ensures
  - Memcached deployment exists
  - Deployment size
- Set default values in `roles/memcached/defaults/main.yml`, in our case: `size: 1`
- Add size in `config/samples/cache_v1alpha1_memcached.yaml`, `size: 2`

## Configure the image regisrty

- Make changes to Makefile to match docker hub settings
- build and push to docker hub
  
  ```sh
  make docker-build docker-push
  ```

## Create the memcached CR

```sh
kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
```

## Deploy the operator

```sh
make deploy
```

## Verify the custom resource is created

- check memcached CR: `kubectl get memcached/memcached-sample -o yaml`
- There shouold be a new namespace similar to `ansible-system`: ``kubectl get ns`
- A new custom resource should be created: `kubectl api-resource | grep memcached`
- `kubectl get deploy -n ansible-system`
