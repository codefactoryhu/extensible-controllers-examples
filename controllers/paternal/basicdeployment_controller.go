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

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	paternalv1 "github.com/codefactoryhu/extensible-controllers/apis/paternal/v1"
)

// BasicDeploymentReconciler reconciles a BasicDeployment object
type BasicDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = paternalv1.GroupVersion.String()
)

//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicdeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicreplicasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicreplicasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paternal.codefactory.hu,resources=basicreplicasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BasicDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *BasicDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var basicDeployment paternalv1.BasicDeployment
	if err := r.Get(ctx, req.NamespacedName, &basicDeployment); err != nil {
		log.Error(err, "unable to fetch BasicDeployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var childReplicaSets paternalv1.BasicReplicaSetList
	if err := r.List(ctx, &childReplicaSets, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child BasicReplicaSets")
		return ctrl.Result{}, err
	}

	var Ready bool

	if len(childReplicaSets.Items) == 1 {
		childReplicaSet := childReplicaSets.Items[0]
		if childReplicaSet.Status.ActiveReplicas == childReplicaSet.Spec.ReplicaCount {
			Ready = true
		} else {
			Ready = false
		}
	} else {
		Ready = false
	}

	basicDeployment.Status.Ready = Ready

	log.V(1).Info("BasicDeployment is ready", "Ready", Ready)
	if err := r.Status().Update(ctx, &basicDeployment); err != nil {
		log.Error(err, "unable to update BasicDeployment status")
		return ctrl.Result{}, err
	}

	if len(childReplicaSets.Items) == 0 {
		basicReplicaSet, err := r.constructBasicReplicaSet(&basicDeployment)
		if err != nil {
			log.Error(err, "unable to construct BasicReplicaSet")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, basicReplicaSet); err != nil {
			log.Error(err, "unable to create BasicReplicaSet")
			return ctrl.Result{}, err
		}
		log.V(1).Info("created Job for CronJob run", "BasicReplicaSet", basicReplicaSet.Name)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BasicDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// if err := mgr.GetFieldIndexer().IndexField(context.Background(), &paternalv1.BasicReplicaSet{}, jobOwnerKey, func(rawObj client.Object) []string {
	// 	// grab the replicaSet object, extract the owner...
	// 	basicReplicaSet := rawObj.(*paternalv1.BasicReplicaSet)
	// 	owner := v1.GetControllerOf(basicReplicaSet)
	// 	if owner == nil {
	// 		return nil
	// 	}
	// 	// ...make sure it's a BasicDeployment...
	// 	if owner.APIVersion != apiGVStr || owner.Kind != "BasicDeployment" {
	// 		return nil
	// 	}

	// 	// ...and if so, return it
	// 	return []string{owner.Name}
	// }); err != nil {
	// 	return err
	// }
	return ctrl.NewControllerManagedBy(mgr).
		Named("basicdeployment").
		For(&paternalv1.BasicDeployment{}).
		Owns(&paternalv1.BasicReplicaSet{}).
		Complete(r)
}

func (r *BasicDeploymentReconciler) constructBasicReplicaSet(basicDeployment *paternalv1.BasicDeployment) (*paternalv1.BasicReplicaSet, error) {
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	name := basicDeployment.Name

	replicaSet := &paternalv1.BasicReplicaSet{
		ObjectMeta: v1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        name,
			Namespace:   basicDeployment.Namespace,
		},
		Spec: paternalv1.BasicReplicaSetSpec{
			ReplicaCount: basicDeployment.Spec.Replicas,
		},
	}
	if err := ctrl.SetControllerReference(basicDeployment, replicaSet, r.Scheme); err != nil {
		return nil, err
	}

	return replicaSet, nil
}
