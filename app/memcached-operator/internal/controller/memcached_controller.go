/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logrun "sigs.k8s.io/controller-runtime/pkg/log"

	cachev1 "github.com/urans/kubemaze/app/memcached-operator/api/v1"
)

const (
	typeStateAvilable = "Available"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.urans.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.urans.com,resources=memcacheds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.urans.com,resources=memcacheds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logrun.FromContext(ctx)
	logger.Info("Reconciling Memcached", "namespace", req.Namespace, "name", req.Name)

	memcached := &cachev1.Memcached{}
	err := r.Get(ctx, req.NamespacedName, memcached)
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Info("Memcached resource not found - ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Memcached resource")
		return ctrl.Result{}, err
	}

	if len(memcached.Status.Conditions) == 0 {
		meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
			Type:    typeStateAvilable,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Memcached is starting reconciliation",
		})
		if err := r.Status().Update(ctx, memcached); err != nil {
			logger.Error(err, "Failed to update Memcached status")
			return ctrl.Result{}, err
		}
		if err := r.Get(ctx, req.NamespacedName, memcached); err != nil {
			logger.Error(err, "Failed to refetch Memcached resource")
			return ctrl.Result{}, err
		}
	}

	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      memcached.Name,
		Namespace: memcached.Namespace,
	}, found)
	if err != nil && apierrs.IsNotFound(err) {
		deploy, err := r.deploymentForMemcached(memcached)
		if err != nil {
			logger.Error(err, "Failed to create deployment for Memcached")

			meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
				Type:    typeStateAvilable,
				Status:  metav1.ConditionFalse,
				Reason:  "Reconciling",
				Message: "Failed to create deployment for Memcached",
			})

			if err := r.Status().Update(ctx, memcached); err != nil {
				logger.Error(err, "Failed to update Memcached status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		logger.Info("Creating a new Deployment for Memcached",
			"namespace", memcached.Namespace, "name", memcached.Name)
		if err := r.Create(ctx, deploy); err != nil {
			logger.Error(err, "Failed to create new Deployment for Memcached",
				"namespace", memcached.Namespace, "name", memcached.Name)
			return ctrl.Result{}, err
		}

		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Memcached deployment")
		return ctrl.Result{}, err
	}

	size := memcached.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		if err := r.Update(ctx, found); err != nil {
			logger.Error(err, "Failed to update Memcached deployment replicas",
				"namespace", memcached.Namespace, "name", memcached.Name)

			if err := r.Get(ctx, req.NamespacedName, memcached); err != nil {
				logger.Error(err, "Failed to refetch Memcached resource after update")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
				Type:    typeStateAvilable,
				Status:  metav1.ConditionFalse,
				Reason:  "Resizing",
				Message: "Failed to update Memcached deployment replicas",
			})

			if err := r.Status().Update(ctx, memcached); err != nil {
				logger.Error(err, "Failed to update Memcached status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeStateAvilable,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciled",
		Message: "Memcached deployment is available with the desired number of replicas",
	})

	if err := r.Status().Update(ctx, memcached); err != nil {
		logger.Error(err, "Failed to update Memcached status after reconciliation")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.Memcached{}).
		Named("memcached").
		Complete(r)
}

// deploymentForMemcached returns a Memcached Deployment object
func (r *MemcachedReconciler) deploymentForMemcached(memcached *cachev1.Memcached) (*appsv1.Deployment, error) {
	replicas := memcached.Spec.Size
	image := "memcached:1.6.26-alpine3.19"

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcached.Name,
			Namespace: memcached.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app.kubernetes.io/name": "project"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app.kubernetes.io/name": "project"},
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: ptr.To(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "memcached",
						ImagePullPolicy: corev1.PullIfNotPresent,
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot:             ptr.To(true),
							RunAsUser:                ptr.To(int64(1001)),
							AllowPrivilegeEscalation: ptr.To(false),
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{
									"ALL",
								},
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
						Command: []string{"memcached", "--memory-limit=64", "-o", "modern", "-v"},
					}},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(memcached, deploy, r.Scheme); err != nil {
		return nil, err
	}
	return deploy, nil
}
