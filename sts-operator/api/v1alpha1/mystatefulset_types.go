/*
Copyright 2024.

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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	LABEL_MYSTS_CONTROLLER_KEY   = "controller-by"
	LABEL_MYSTS_CONTROLLER_VALUE = "myStateFulSetController"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyStatefulSetSpec defines the desired state of MyStatefulSet
type MyStatefulSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Replicas int                `json:"replicas,omitempty"`
	Template v1.PodTemplateSpec `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// MyStatefulSetStatus defines the observed state of MyStatefulSet
type MyStatefulSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodStatus []v1.PodStatus `json:"pod,omitempty"`

	Phase string `json:"phase,omitempty"`

	//Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

type PodStatus struct {
	Name  string `json:"name,omitempty"`
	Phase string `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyStatefulSet is the Schema for the mystatefulsets API
type MyStatefulSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyStatefulSetSpec   `json:"spec,omitempty"`
	Status MyStatefulSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyStatefulSetList contains a list of MyStatefulSet
type MyStatefulSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyStatefulSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyStatefulSet{}, &MyStatefulSetList{})
}
