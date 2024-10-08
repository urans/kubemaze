package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/brianvoe/gofakeit/v6"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubeClient creates a kubernetes client
func NewKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		slog.Error("load kubeconfig failed", "error", err)
		return nil, fmt.Errorf("NewKubeClient failed: %w", err)
	}
	slog.Debug("kubeconfig loaded", "config", config)
	return kubernetes.NewForConfig(config)
}

func createConfigMap(clientset kubernetes.Interface, name, ns string) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"User":  gofakeit.Username(),
			"Phone": gofakeit.Phone(),
		},
	}

	val, err := clientset.CoreV1().ConfigMaps(ns).Create(context.TODO(), cm, metav1.CreateOptions{})
	if err != nil {
		slog.Error("create configmap failed", "err", err)
		return err
	}
	slog.Info("configmap created", "name", val.Name)
	return nil
}

func main() {
	conf := path.Join(os.Getenv("HOME"), ".kube/config-dev")
	client, err := NewKubeClient(conf)
	if err != nil {
		os.Exit(1)
	}

	name := "trace-create-cm-" + gofakeit.Phone()
	err = createConfigMap(client, name, "kallen")
	if err != nil {
		os.Exit(1)
	}
}
