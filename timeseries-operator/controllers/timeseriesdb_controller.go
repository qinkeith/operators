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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "qinkeith.com/operators/timeseries-operator/api/v1"
)

// TimeseriesDBReconciler reconciles a TimeseriesDB object
type TimeseriesDBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operator.qinkeith.com,resources=timeseriesdbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.qinkeith.com,resources=timeseriesdbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.qinkeith.com,resources=timeseriesdbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TimeseriesDB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *TimeseriesDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := log.FromContext(ctx)

	log = log.WithValues("timeseriesdb", req.NamespacedName)

	timeseriesdb := new(operatorv1.TimeseriesDB)

	if err := r.Client.Get(ctx, req.NamespacedName, timeseriesdb); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log = log.WithValues("dbType", timeseriesdb.Spec.DBType, "replicas", timeseriesdb.Spec.Replicas)

	if timeseriesdb.Status.Status == "" || timeseriesdb.Status.Message == "" {
		timeseriesdb.Status = operatorv1.TimeseriesDBStatus{Status: "Initialized", Message: "Database creation is in progress"}
		err := r.Status().Update(ctx, timeseriesdb)
		if err != nil {
			log.Error(err, "status update failed")
			return ctrl.Result{}, err
		}
		log.Info("status updated")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TimeseriesDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.TimeseriesDB{}).
		Complete(r)
}
