/*
Copyright 2022 Zoltan Magyar.

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

package paternal

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	paternalv1 "github.com/codefactoryhu/extensible-controllers/apis/paternal/v1"
)

// BasicReplicaSetReconciler reconciles a BasicReplicaSet object
type BasicReplicaSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicreplicasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicreplicasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicreplicasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BasicReplicaSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *BasicReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var basicReplicaSet paternalv1.BasicReplicaSet
	if err := r.Get(ctx, req.NamespacedName, &basicReplicaSet); err != nil {
		log.Error(err, "unable to fetch BasicReplicaSet")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	basicReplicaSet.Status.ActiveReplicas = basicReplicaSet.Spec.ReplicaCount

	if err := r.Status().Update(ctx, &basicReplicaSet); err != nil {
		log.Error(err, "unable to update BasicReplicaSet status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BasicReplicaSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("basicreplicaset").
		For(&paternalv1.BasicReplicaSet{}).
		Complete(r)
}
