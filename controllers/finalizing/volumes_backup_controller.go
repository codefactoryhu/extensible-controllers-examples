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

package finalizing

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	finalizingv1 "github.com/codefactoryhu/extensible-controllers/apis/finalizing/v1"
	"github.com/go-logr/logr"
)

// VolumesBackupReconciler finalizes a Machine object
type VolumesBackupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const machineFinalizer = "machine.finalizing.codefactory.hu"

//+kubebuilder:rbac:groups=finalizing.codefactory.hu,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=finalizing.codefactory.hu,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=finalizing.codefactory.hu,resources=machines/finalizers,verbs=update

func (r *VolumesBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Finalizing Machine")

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

	isMachineMarkedToBeDeleted := machine.GetDeletionTimestamp() != nil
	if isMachineMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(machine, machineFinalizer) {
			if err := r.finalizeMachine(reqLogger, machine); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(machine, machineFinalizer)
			err = r.Update(ctx, machine)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(machine, machineFinalizer) {
		controllerutil.AddFinalizer(machine, machineFinalizer)
		if err := r.Update(ctx, machine); err != nil {
			reqLogger.Error(err, "Failed to add finalizer for Machine resource.")
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VolumesBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("volumesbackup").
		For(&finalizingv1.Machine{}).
		Complete(r)
}

func (r *VolumesBackupReconciler) finalizeMachine(reqLogger logr.Logger, m *finalizingv1.Machine) error {
	reqLogger.Info("Finalizing Machine")
	// Do some finalizing stuff
	reqLogger.Info("Successfully finalized Machine")
	return nil
}
