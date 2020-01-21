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

package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	linodego "github.com/linode/linodego"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane-runtime/pkg/meta"
	"github.com/crossplaneio/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplaneio/crossplane-runtime/pkg/resource"

	linodev1alpha1 "github.com/displague/stack-linode/api/v1alpha1"
	"github.com/displague/stack-linode/clients"
)

const (
	errNewClient      = "cannot create new Instance client"
	errNotInstance    = "managed resource is not an Instance"
	errInstanceCreate = "cannot create Instance"
	errInstanceDelete = "cannot delete Instance"
)

// InstanceController is responsible for adding the Instance
// controller and its corresponding reconciler to the manager with any runtime configuration.
type InstanceController struct{}

var (
	controllerLog = ctrl.Log.WithName("instance.controller")
)

// SetupWithManager creates a new Instance Controller and adds it to the
// Manager with default RBAC. The Manager will set fields on the Controller and
// start it when the Manager is Started.
func (c *InstanceController) SetupWithManager(mgr ctrl.Manager) error {
	r := managed.NewReconciler(mgr,
		resource.ManagedKind(linodev1alpha1.InstanceGroupVersionKind),
		managed.WithConnectionPublishers(),
		managed.WithExternalConnecter(&connecter{client: mgr.GetClient()}))

	name := strings.ToLower(fmt.Sprintf("%s.%s", linodev1alpha1.InstanceKind, linodev1alpha1.Group))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&linodev1alpha1.Instance{}).
		Complete(r)
}

type connecter struct {
	client      client.Client
	newClientFn func(credentials []byte) linodego.Client
}

// Connect to the supplied resource.Managed (presumed to be an
// Instance) by using the Provider it references to create a new
// Linode API client.
func (c *connecter) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		err := errors.New(errNotInstance)
		controllerLog.Error(err, "Connect", "mg", mg)
		return nil, err
	}

	controllerLog.Info("Connect", "spec", m.Spec, "status", m.Status)

	p := &linodev1alpha1.Provider{}
	n := meta.NamespacedNameOf(m.Spec.ProviderReference)
	if err := c.client.Get(ctx, n, p); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider %s", n)
	}

	s := &corev1.Secret{}
	n = types.NamespacedName{Namespace: p.GetNamespace(), Name: p.Spec.Secret.Name}
	if err := c.client.Get(ctx, n, s); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider secret %s", n)
	}
	newClientFn := clients.NewClient
	if c.newClientFn != nil {
		newClientFn = c.newClientFn
	}
	client := newClientFn(s.Data[p.Spec.Secret.Key])
	return &external{client: client}, errors.Wrap(nil, errNewClient)
}

type external struct{ client linodego.Client }

// Observe the existing external resource, if any. The resource.ManagedReconciler
// calls Observe in order to determine whether an external resource needs to be
// created, updated, or deleted.
func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotInstance)
	}

	controllerLog.Info("Observe", "spec", m.Spec, "status", m.Status)

	if m.Status.AtProvider.Id == 0 {
		return managed.ExternalObservation{}, nil
	}

	instance, err := e.client.GetInstance(ctx, m.Status.AtProvider.Id)

	controllerLog.Info("Observe", "instanceId", m.Status.AtProvider.Id, "err", err)

	if err != nil {
		if e, ok := err.(*linodego.Error); ok && e.Code == http.StatusNotFound {
			return managed.ExternalObservation{}, nil
		}
	}

	controllerLog.Info("Observe", "wantLabel", m.Spec.ForProvider.Label, "gotLabel", instance.Label)
	switch m.Status.AtProvider.Status {
	case string(linodego.InstanceRunning):
		m.Status.SetConditions(runtimev1alpha1.Available())
		resource.SetBindable(m)
	case string(linodego.InstanceProvisioning):
		m.Status.SetConditions(runtimev1alpha1.Creating())
	}

	// Store observed values in Status
	m.Status.AtProvider.Id = instance.ID
	m.Status.AtProvider.Label = instance.Label
	m.Status.AtProvider.Status = string(instance.Status)
	m.Status.AtProvider.Region = instance.Region
	m.Status.AtProvider.Type = instance.Type
	m.Status.AtProvider.Image = instance.Image
	m.Status.AtProvider.IPv4 = []string{}
	for _, ip := range instance.IPv4 {
		m.Status.AtProvider.IPv4 = append(m.Status.AtProvider.IPv4, ip.String())
	}
	m.Status.AtProvider.IPv6 = instance.IPv6

	// Compare observed (GetInstance()) to desired (spec)
	upToDate := m.Spec.ForProvider.Label == "" || instance.Label == m.Spec.ForProvider.Label
	isOnOrOff := map[string]bool{
		string(linodego.InstanceRunning): true,
		string(linodego.InstanceOffline): true,
	}

	needsPowerToggle := (!isOnOrOff[string(instance.Status)] || m.Spec.ForProvider.Status != string(instance.Status))
	upToDate = upToDate && !needsPowerToggle

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: upToDate,
	}, nil
}

// Create a new external resource based on the specification of our managed
// resource. resource.ManagedReconciler only calls Create if Observe reported
// that the external resource did not exist.
func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotInstance)
	}
	controllerLog.Info("Create", "spec", m.Spec, "status", m.Status)

	m.Status.SetConditions(runtimev1alpha1.Creating())

	booted := m.Spec.ForProvider.Status == string(linodego.InstanceRunning)
	rootPass, _ := createRandomRootPassword()
	instance, err := e.client.CreateInstance(ctx, linodego.InstanceCreateOptions{
		Label:           m.Spec.ForProvider.Label,
		Region:          m.Spec.ForProvider.Region,
		Type:            m.Spec.ForProvider.Type,
		AuthorizedUsers: m.Spec.ForProvider.AuthorizedUsers,
		Image:           m.Spec.ForProvider.Image,
		Booted:          &booted,
		RootPass:        rootPass,
	})
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errInstanceCreate)
	}
	m.Status.SetConditions(runtimev1alpha1.Available())

	m.Status.AtProvider.Id = instance.ID

	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{
			"rootPass": []byte(rootPass),
			"ipv6":     []byte(instance.IPv6),
		},
	}, nil
}

// Update the existing external resource to match the specifications of our
// managed resource. resource.ManagedReconciler only calls Update if Observe
// reported that the external resource was not up to date.
func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	var err error
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotInstance)
	}

	instance, errGetting := e.client.GetInstance(ctx, m.Status.AtProvider.Id)
	if errGetting != nil {
		return managed.ExternalUpdate{}, err
	}

	if m.Spec.ForProvider.Status == string(linodego.InstanceOffline) &&
		instance.Status == linodego.InstanceRunning {
		err = e.client.ShutdownInstance(ctx, m.Status.AtProvider.Id)
	} else if instance.Status == linodego.InstanceOffline {
		err = e.client.BootInstance(ctx, m.Status.AtProvider.Id, 0)
	}
	controllerLog.Info("Update", "spec", m.Spec, "status", m.Status)

	return managed.ExternalUpdate{}, err
}

// Delete the external resource. resource.ManagedReconciler only calls Delete
// when a managed resource with the 'Delete' reclaim policy has been deleted.
func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return errors.New(errNotInstance)
	}
	controllerLog.Info("Delete", "spec", m.Spec, "status", m.Status)

	m.SetConditions(runtimev1alpha1.Deleting())
	err := e.client.DeleteInstance(ctx, m.Status.AtProvider.Id)

	if err != nil {
		if e, ok := err.(*linodego.Error); ok && e.Code == http.StatusNotFound {
			return nil
		}
	}

	return errors.Wrap(err, errInstanceDelete)
}

func createRandomRootPassword() (string, error) {
	rawRootPass := make([]byte, 50)
	_, err := rand.Read(rawRootPass)
	if err != nil {
		return "", fmt.Errorf("Failed to generate random password")
	}
	rootPass := base64.StdEncoding.EncodeToString(rawRootPass)
	return rootPass, nil
}
