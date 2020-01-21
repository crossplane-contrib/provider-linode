/*
Copyright 2019 The Crossplane Authors.

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
	"reflect"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

var (
	LkeKind             = reflect.TypeOf(Lke{}).Name()
	LkeKindAPIVersion   = LkeKind + "." + GroupVersion.String()
	LkeGroupVersionKind = GroupVersion.WithKind(LkeKind)
)

// +kubebuilder:validation:Required

// LkeParameters defines the desired state of Lke
type LkeParameters struct {
	NodePools []LkeClusterPool `json:"node_pools"`
	Label     string           `json:"label"`
	Region    string           `json:"region"`
	Version   string           `json:"version"`
	Tags      []string         `json:"tags,omitempty"`
}

// LkeObservation defines the observed state of Lke
type LkeObservation struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Id      int      `json:"id"`
	Created string   `json:"created"`
	Updated string   `json:"updated"`
	Label   string   `json:"label"`
	Region  string   `json:"region"`
	Status  string   `json:"status"`
	Version string   `json:"version"`
	Tags    []string `json:"tags"`

	NodePools []LkeClusterPool `json:"node_pools"`
}

// LkeStatus defines the observed state of Lke
type LkeStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     LkeObservation `json:"atProvider"`
}

// LkeSpec defines the desired state of LKE
type LkeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  LkeParameters `json:"forProvider"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Lke is the Schema for the lkes API
type Lke struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec LkeSpec `json:"spec,omitempty"`

	// +optional
	Status LkeStatus `json:"status,omitempty"`
}

// LkeClusterPoolLinode
type LkeClusterPoolLinode struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

// LkeClusterPool
type LkeClusterPool struct {
	Id      int                    `json:"id,omitempty"`
	Count   int                    `json:"count,omitempty"`
	Type    string                 `json:"type,omitempty"`
	Linodes []LkeClusterPoolLinode `json:"linodes,omitempty"`
}

// +kubebuilder:object:root=true

// LkeList contains a list of Lke
type LkeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Lke `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Lke{}, &LkeList{})
}
