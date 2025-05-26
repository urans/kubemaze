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

- Build Image

```bash
❯ make docker-build docker-push IMG=kelein/memcached-operator:v0.1.1
```

- Install CRD

```bash
❯ make install
```

- Install Controller

```bash
❯ make deploy
# OR
❯ make deploy IMG=kelein/memcached-operator:v0.1.1
```
