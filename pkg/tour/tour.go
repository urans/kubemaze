package tour

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
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
