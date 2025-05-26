# Memcached Operator

- Create project

```bash
❯ kubebuilder init --domain urans.com --repo github.com/urans/kubemaze/app/memcached-operator
```

- Create API

```bash
❯ kubebuilder create api --group cache --version v1alpha1 --kind Memcached
```

```bash
❯ make generate
```
