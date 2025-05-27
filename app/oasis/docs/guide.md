# Kubebuilder Guide

```bash
❯ kubebuilder init --domain oasis.urans.com --repo github.com/urans/kubemaze/app/oasis

❯ kubebuilder create api --group webapp --version v1 --kind Oasis



❯ make manifests

❯ make install

❯ kubectl apply -k config/samples/

❯ make docker-build docker-push IMG=kelein/kubemaze/guestbook:v0.1.0

make deploy IMG=<some-registry>/<project-name>:tag

```
