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
