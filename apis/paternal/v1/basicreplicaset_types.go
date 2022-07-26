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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BasicReplicaSetSpec defines the desired state of BasicReplicaSet
type BasicReplicaSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of BasicReplicaSet. Edit basicreplicaset_types.go to remove/update
	ReplicaCount *int64 `json:"replicaCount"`
}

// BasicReplicaSetStatus defines the observed state of BasicReplicaSet
type BasicReplicaSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+optional
	ActiveReplicas *int64 `json:"activeReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BasicReplicaSet is the Schema for the basicreplicasets API
type BasicReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BasicReplicaSetSpec   `json:"spec,omitempty"`
	Status BasicReplicaSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BasicReplicaSetList contains a list of BasicReplicaSet
type BasicReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BasicReplicaSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BasicReplicaSet{}, &BasicReplicaSetList{})
}
