package main

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/urans/kubemaze/pkg/tour"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func initLogger() {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	// * Text Log Format
	logger := slog.New(slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			ReplaceAttr: replace,
		},
	))
	slog.SetDefault(logger)
}

func main() {
	initLogger()

	local := path.Join(os.Getenv("HOME"), ".kube/config")
	clientset, err := tour.NewKubeClient(local)
	if err != nil {
		slog.Error("create kube client failed", "error", err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Minute)
	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			val := obj.(v1.Object)
			slog.Info("New Pod Created", "namespace", val.GetNamespace(), "pod", val.GetName())
		},
		UpdateFunc: func(oldObj, newObj any) {
			old := oldObj.(v1.Object)
			another := newObj.(v1.Object)
			slog.Info("Pod Updated", "namespace", old.GetNamespace(), "old", old.GetName(), "new", another.GetName())
		},
		DeleteFunc: func(obj any) {
			val := obj.(v1.Object)
			slog.Info("Pod Deleted", "namespace", val.GetNamespace(), "pod", val.GetName())
		},
	})

	quit := make(chan struct{})
	defer close(quit)
	informer.Run(quit)
}
