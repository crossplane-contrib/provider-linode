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
	"fmt"
	"strings"

	linodego "github.com/linode/linodego"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane-runtime/pkg/meta"
	"github.com/crossplaneio/crossplane-runtime/pkg/resource"

	linodev1alpha1 "github.com/displague/stack-linode/api/v1alpha1"
	"github.com/displague/stack-linode/clients"
)

const (
	errNewClient      = "cannot create new Instance client"
	errNotInstance    = ""
	errInstanceCreate = ""
)

// InstanceController is responsible for adding the Instance
// controller and its corresponding reconciler to the manager with any runtime configuration.
type InstanceController struct{}

// SetupWithManager creates a new Instance Controller and adds it to the
// Manager with default RBAC. The Manager will set fields on the Controller and
// start it when the Manager is Started.
func (c *InstanceController) SetupWithManager(mgr ctrl.Manager) error {
	r := resource.NewManagedReconciler(mgr,
		resource.ManagedKind(linodev1alpha1.InstanceGroupVersionKind),
		resource.WithManagedConnectionPublishers(),
		resource.WithExternalConnecter(&connecter{client: mgr.GetClient()}))

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

func (c *connecter) Connect(ctx context.Context, mg resource.Managed) (resource.ExternalClient, error) {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return nil, errors.New(errNotInstance)
	}

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

func (e *external) Observe(ctx context.Context, mg resource.Managed) (resource.ExternalObservation, error) {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return resource.ExternalObservation{}, errors.New(errNotInstance)
	}

	if m.Status.Id == 0 {
		return resource.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	return resource.ExternalObservation{ResourceExists: true}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (resource.ExternalCreation, error) {
	m, ok := mg.(*linodev1alpha1.Instance)
	if !ok {
		return resource.ExternalCreation{}, errors.New(errNotInstance)
	}

	m.Status.SetConditions(runtimev1alpha1.Creating())

	instance, err := e.client.CreateInstance(ctx, linodego.InstanceCreateOptions{
		Label: m.Spec.Label,
	})
	if err != nil {
		return resource.ExternalCreation{}, errors.Wrap(err, errInstanceCreate)
	}

	m.Status.Id = instance.ID

	return resource.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (resource.ExternalUpdate, error) {
	// m, ok := mg.(*linodev1alpha1.Instance)
	// if !ok {
	// 	return resource.ExternalUpdate{}, errors.New(errNotInstance)
	// }

	return resource.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	// m, ok := mg.(*linodev1alpha1.Instance)
	// if !ok {
	// 	return errors.New(errNotInstance)
	// }

	return nil
}
