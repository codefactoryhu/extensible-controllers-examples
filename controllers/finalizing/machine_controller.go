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

package finalizing

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	finalizingv1 "github.com/codefactoryhu/extensible-controllers/apis/finalizing/v1"
)

// MachineReconciler reconciles a Machine object
type MachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=finalizing.codefactory.hu,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=finalizing.codefactory.hu,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=finalizing.codefactory.hu,resources=machines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Machine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *MachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Reconciling Machine")

	machine := &finalizingv1.Machine{}
	err := r.Get(ctx, req.NamespacedName, machine)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Machine resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		reqLogger.Error(err, "Failed to get Machine resource.")
		return ctrl.Result{}, err
	}

	// Do reconciliation stuff

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("machine").
		For(&finalizingv1.Machine{}).
		Complete(r)
}
