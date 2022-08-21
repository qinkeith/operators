/*
Copyright 2022.

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

package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/docker/docker"

	corev1alpha1 "github.com/qinkeith/operators/clair-opertor/api/v1alpha1"
)

// ScannerReconciler reconciles a Scanner object
type ScannerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.qinkeith.com,resources=scanners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.qinkeith.com,resources=scanners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.qinkeith.com,resources=scanners/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scanner object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *ScannerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var log = log.FromContext(ctx)

	var (
		scanners corev1alpha1.ScannerList
		pod      corev1.Pod
	)

	opts := []client.ListOption{
		client.InNamespace(req.NamespacedName.Namespace),
	}
	if err := r.List(ctx, &scanners, opts...); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(scanners.Items) == 0 {
		log.Info("No Scanners found for namespace " + req.NamespacedName.Namespace)
		return ctrl.Result{}, nil
	}

	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		log.Error(err, "unable to fetch Pod")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	podFinalizer := "qinkeith.com/podFinalizer"

	isPodMarkedToBeDeleted := pod.GetDeletionTimestamp() != nil
	if isPodMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(&pod, podFinalizer) {
			if err := r.finalizePod(ctx, &pod, scanners.Items[0].Spec.ClairBaseUrl, scanners.Items[0].Spec.SlackWebhookUrl); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(&pod, podFinalizer)
			_ = r.Update(ctx, &pod)
		}
	}
	if !controllerutil.ContainsFinalizer(&pod, podFinalizer) {
		controllerutil.AddFinalizer(&pod, podFinalizer)
		_ = r.Update(ctx, &pod)
	}

	return ctrl.Result{}, nil
}

func (r *ScannerReconciler) finalizePod(ctx context.Context, pod *corev1.Pod, clairBaseUrl string, slackWebhookUrl string) error {
	for _, container := range pod.Status.ContainerStatuses {
		manifest, err := docker.Inspect(ctx, container.Image)
		if err != nil {
			log.Error(err, "Error while inspecting container Manifest")
			continue
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScannerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Scanner{}).
		Watches(&source.Kind{Type: &corev1.Pod{}},
			handler.EnqueueRequestsFromMapFunc(
				func(pod client.Object) []reconcile.Request {
					return []reconcile.Request{
						{NamespacedName: types.NamespacedName{
							Name:      pod.GetName(),
							Namespace: pod.GetNamespace(),
						}},
					}
				},
			),
		).
		Complete(r)
}
