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

// BasicStatefulSetSpec defines the desired state of BasicStatefulSet
type BasicStatefulSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of BasicStatefulSet. Edit basicstatefulset_types.go to remove/update
	ReplicaCount *int `json:"replicaCount"`
}

// BasicStatefulSetStatus defines the observed state of BasicStatefulSet
type BasicStatefulSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+optional
	ActiveReplicas *int `json:"activeReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BasicStatefulSet is the Schema for the basicstatefulsets API
type BasicStatefulSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BasicStatefulSetSpec   `json:"spec,omitempty"`
	Status BasicStatefulSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BasicStatefulSetList contains a list of BasicStatefulSet
type BasicStatefulSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BasicStatefulSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BasicStatefulSet{}, &BasicStatefulSetList{})
}
