package main

import (
	"context"
	"log/slog"
	"os"
	"path"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/urans/kubemaze/pkg/tour"
)

var watchTimeoutSeconds = int64(10)

func main() {
	clientset, err := tour.NewKubeClient(path.Join(os.Getenv("HOME"), ".kube/config"))
	if err != nil {
		slog.Error("create kube client failed", "error", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		watchNamespaces(clientset)
	}()
	wg.Wait()
}

func watchNamespaces(clientset kubernetes.Interface) error {
	watcher, err := clientset.CoreV1().Namespaces().Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &watchTimeoutSeconds})
	if err != nil {
		return err
	}

	for event := range watcher.ResultChan() {
		ns := event.Object.(*corev1.Namespace)
		switch event.Type {
		case watch.Modified:
		case watch.Bookmark:
		case watch.Deleted:
		case watch.Error:
		case watch.Added:
			processNamespaces(ns.GetName(), event.Type)
		}
	}
	return nil
}

func processNamespaces(namespace string, event watch.EventType) {
	slog.Info("received event", "type", event, "namespaces", namespace)
}
