package tour

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const masterLabelName = "node-role.kubernetes.io/master"

// NewKubeClient creates a kubernetes client
func NewKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		slog.Error("load kubeconfig failed", "error", err)
		return nil, err
	}
	slog.Debug("kubeconfig loaded", "config", config)
	return kubernetes.NewForConfig(config)
}

// ListNodes list all nodes in the cluster
func ListNodes(clientset kubernetes.Interface, labels map[string]string) ([]corev1.Node, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: buildLabelSelector(labels),
	})
	if err != nil {
		slog.Error("list nodes failed", "error", err)
		return nil, err
	}
	return nodes.Items, nil
}

func buildLabelSelector(labels map[string]string) string {
	pairs := []string{}
	for k, v := range labels {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, ",")
}

// GetNode get node detail in the cluster
func GetNode(clientset kubernetes.Interface, name string) (*corev1.Node, error) {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		slog.Error("get node failed", "name", name, "error", err)
		return nil, err
	}
	return node, nil
}

func isMaster(node *corev1.Node) bool {
	if node == nil {
		return false
	}
	_, ok := node.Labels[masterLabelName]
	return ok
}

func isReady(node *corev1.Node) bool {
	if node == nil {
		return false
	}
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func nodeAge(node *corev1.Node) time.Duration {
	if node == nil {
		return 0
	}
	return time.Now().Sub(node.CreationTimestamp.Time)
}

func kubeletVersion(node *corev1.Node) string {
	if node == nil {
		return ""
	}
	return node.Status.NodeInfo.KubeletVersion
}

// ListPods list all pod in specified namespace
func ListPods(clientset kubernetes.Interface, namespace string) ([]corev1.Pod, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		slog.Error("list pods failed", "error", err)
		return nil, err
	}
	return pods.Items, nil
}

// GetPod get pod detail in specified namespace and name
func GetPod(clientset kubernetes.Interface, namespace, name string) (*corev1.Pod, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		slog.Error("get pod failed", "name", name, "error", err)
		return nil, err
	}
	return pod, nil
}

// CreateSecretFromFile creates a secret from a file
func CreateSecretFromFile(clientset kubernetes.Interface, namespace, name, fpath string) (*corev1.Secret, error) {
	content, err := os.ReadFile(fpath)
	if err != nil {
		slog.Error("read file failed", "path", fpath, "error", err)
		return nil, err
	}
	dir, fname := path.Split(fpath)
	slog.Info("file info", "dir", dir, "name", fname)
	data := make(map[string][]byte)
	data[fname] = content

	result, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(),
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: name, Namespace: namespace,
			},
			Data: data,
		}, metav1.CreateOptions{},
	)
	if kerrs.IsAlreadyExists(err) {
		slog.Warn(err.Error(), "namespace", namespace, "name", name)
		return result, nil
	}
	return result, err
}
