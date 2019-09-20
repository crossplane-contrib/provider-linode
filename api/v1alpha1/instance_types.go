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
	"net"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	Group   = "compute.linode.crossplane.io"
	Version = "v1alpha1"
)

var (

	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	InstanceKind             = reflect.TypeOf(Instance{}).Name()
	InstanceKindAPIVersion   = InstanceKind + "." + SchemeGroupVersion.String()
	InstanceGroupVersionKind = SchemeGroupVersion.WithKind(InstanceKind)
)

type InstanceParameters struct {
	// +optional
	Label string `json:",omitempty"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`

	InstanceParameters string `json:",inline"`
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ID     int      `json:"id"`
	IPv4   []net.IP `json:"ipv4"`
	IPv6   string   `json:"ipv6"`
	Label  string   `json:"label"`
	Status string   `json:"status"`
	Region string   `json:"region"`
	Type   string   `json:"type"`
}

// +kubebuilder:object:root=true

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// SetBindingPhase of this Instance.
func (a *Instance) SetBindingPhase(p runtimev1alpha1.BindingPhase) {
	a.Status.SetBindingPhase(p)
}

// GetBindingPhase of this Instance.
func (a *Instance) GetBindingPhase() runtimev1alpha1.BindingPhase {
	return a.Status.GetBindingPhase()
}

// SetConditions of this Instance.
func (a *Instance) SetConditions(c ...runtimev1alpha1.Condition) {
	a.Status.SetConditions(c...)
}

// SetClaimReference of this Instance.
func (a *Instance) SetClaimReference(r *corev1.ObjectReference) {
	a.Spec.ClaimReference = r
}

// GetClaimReference of this Instance.
func (a *Instance) GetClaimReference() *corev1.ObjectReference {
	return a.Spec.ClaimReference
}

// SetNonPortableClassReference of this Instance.
func (i *Instance) SetNonPortableClassReference(r *corev1.ObjectReference) {
	i.Spec.NonPortableClassReference = r
}

// GetNonPortableClassReference of this Instance.
func (i *Instance) GetNonPortableClassReference() *corev1.ObjectReference {
	return i.Spec.NonPortableClassReference
}

// SetWriteConnectionSecretToReference of this Instance.
func (a *Instance) SetWriteConnectionSecretToReference(r corev1.LocalObjectReference) {
	a.Spec.WriteConnectionSecretToReference = r
}

// GetWriteConnectionSecretToReference of this Instance.
func (a *Instance) GetWriteConnectionSecretToReference() corev1.LocalObjectReference {
	return a.Spec.WriteConnectionSecretToReference
}

// GetReclaimPolicy of this Instance.
func (a *Instance) GetReclaimPolicy() runtimev1alpha1.ReclaimPolicy {
	return a.Spec.ReclaimPolicy
}

// SetReclaimPolicy of this Instance.
func (a *Instance) SetReclaimPolicy(p runtimev1alpha1.ReclaimPolicy) {
	a.Spec.ReclaimPolicy = p
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
