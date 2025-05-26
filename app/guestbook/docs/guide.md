# Kubebuilder Guide

```bash
❯ kubebuilder init --domain urans.com --repo github.com/urans/kubemaze/app/guestbook

❯ kubebuilder create api --group webapp --version v1 --kind Guestbook

❯ make manifests

❯ make install

❯ kubectl apply -k config/samples/

❯ make docker-build docker-push IMG=kelein/kubemaze/guestbook:v0.1.0

make deploy IMG=<some-registry>/<project-name>:tag

```
