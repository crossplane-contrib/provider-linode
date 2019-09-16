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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ObjectStorageSpec defines the desired state of ObjectStorage
type ObjectStorageSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// ObjectStorageStatus defines the observed state of ObjectStorage
type ObjectStorageStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ObjectStorage is the Schema for the objectstorages API
type ObjectStorage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObjectStorageSpec   `json:"spec,omitempty"`
	Status ObjectStorageStatus `json:"status,omitempty"`
}

// SetBindingPhase of this ObjectStorage.
func (a *ObjectStorage) SetBindingPhase(p runtimev1alpha1.BindingPhase) {
	a.Status.SetBindingPhase(p)
}

// GetBindingPhase of this ObjectStorage.
func (a *ObjectStorage) GetBindingPhase() runtimev1alpha1.BindingPhase {
	return a.Status.GetBindingPhase()
}

// SetConditions of this ObjectStorage.
func (a *ObjectStorage) SetConditions(c ...runtimev1alpha1.Condition) {
	a.Status.SetConditions(c...)
}

// SetClaimReference of this ObjectStorage.
func (a *ObjectStorage) SetClaimReference(r *corev1.ObjectReference) {
	a.Spec.ClaimReference = r
}

// GetClaimReference of this ObjectStorage.
func (a *ObjectStorage) GetClaimReference() *corev1.ObjectReference {
	return a.Spec.ClaimReference
}

// SetClassReference of this ObjectStorage.
func (a *ObjectStorage) SetClassReference(r *corev1.ObjectReference) {
	a.Spec.ClassReference = r
}

// GetClassReference of this ObjectStorage.
func (a *ObjectStorage) GetClassReference() *corev1.ObjectReference {
	return a.Spec.ClassReference
}

// SetWriteConnectionSecretToReference of this ObjectStorage.
func (a *ObjectStorage) SetWriteConnectionSecretToReference(r corev1.LocalObjectReference) {
	a.Spec.WriteConnectionSecretToReference = r
}

// GetWriteConnectionSecretToReference of this ObjectStorage.
func (a *ObjectStorage) GetWriteConnectionSecretToReference() corev1.LocalObjectReference {
	return a.Spec.WriteConnectionSecretToReference
}

// GetReclaimPolicy of this ObjectStorage.
func (a *ObjectStorage) GetReclaimPolicy() runtimev1alpha1.ReclaimPolicy {
	return a.Spec.ReclaimPolicy
}

// SetReclaimPolicy of this ObjectStorage.
func (a *ObjectStorage) SetReclaimPolicy(p runtimev1alpha1.ReclaimPolicy) {
	a.Spec.ReclaimPolicy = p
}

// +kubebuilder:object:root=true

// ObjectStorageList contains a list of ObjectStorage
type ObjectStorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObjectStorage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ObjectStorage{}, &ObjectStorageList{})
}
