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


