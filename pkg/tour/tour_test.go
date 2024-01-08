package tour

import (
	"log/slog"
	"testing"

	"k8s.io/client-go/kubernetes"
)

const devKubeConfig = "/Users/kallen/.kube/config"

func TestNewKubeClient(t *testing.T) {
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"A", args{"~/.kube/config"}, true},
		{"B", args{devKubeConfig}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKubeClient(tt.args.kubeconfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubeClient() got = %v, error = %v, wantErr %v", got, err, tt.wantErr)
				return
			}
		})
	}
}

func initDevKubeClient(t *testing.T) kubernetes.Interface {
	t.Helper()
	client, _ := NewKubeClient(devKubeConfig)
	return client
}

func TestListNodes(t *testing.T) {
	tests := []struct {
		name    string
		client  kubernetes.Interface
		labels  map[string]string
		wantErr bool
	}{
		{"A", initDevKubeClient(t), nil, false},
		{"B", initDevKubeClient(t), map[string]string{"node-role.kubernetes.io/master": "yes"}, false},
		{"C", initDevKubeClient(t), map[string]string{"minikube": "yes"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListNodes(tt.client, tt.labels)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, node := range got {
				slog.Info("ranged node", "namespace", node.Namespace, "name", node.Name, "kind", node.Kind, "labels", node.Labels)
			}
		})
	}
}

func TestGetNode(t *testing.T) {
	type args struct {
		clientset kubernetes.Interface
		name      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"A", args{initDevKubeClient(t), "minikube"}, false},
		{"B", args{initDevKubeClient(t), ""}, true},
		{"C", args{initDevKubeClient(t), "dev"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNode(tt.args.clientset, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				slog.Info("node info", "name", got.Name, "isMaster", isMaster(got), "isReady", isReady(got), "age", nodeAge(got), "kubelet", kubeletVersion(got))
			}
		})
	}
}
