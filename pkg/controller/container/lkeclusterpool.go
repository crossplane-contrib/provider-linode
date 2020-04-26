/*
Copyright 2020 The Crossplane Authors.
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

package container

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/displague/crossplane-provider-linode/apis/container/v1alpha1"
	apisv1alpha1 "github.com/displague/crossplane-provider-linode/apis/v1alpha1"
)

const (
	errNotLKEClusterPool = "managed resource is not a LKEClusterPool custom resource"
)

// SetupLKEClusterPool adds a controller that reconciles LKEClusterPool managed resources.
func SetupLKEClusterPool(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(v1alpha1.LKEClusterPoolGroupKind)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.LKEClusterPoolGroupVersionKind),
		managed.WithExternalConnecter(&lkeClusterPoolConnector{kube: mgr.GetClient(), newServiceFn: newNoOpService}),
		managed.WithInitializers(managed.NewNameAsExternalName(mgr.GetClient())),
		managed.WithLogger(l.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.LKEClusterPool{}).
		Complete(r)
}

type lkeClusterPoolConnector struct {
	kube         client.Client
	newServiceFn func() (interface{}, error)
}

func (c *lkeClusterPoolConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.LKEClusterPool)
	if !ok {
		return nil, errors.New(errNotLKEClusterPool)
	}

	provider := &apisv1alpha1.Provider{}
	if err := c.kube.Get(ctx, meta.NamespacedNameOf(cr.Spec.ProviderReference), provider); err != nil {
		return nil, errors.Wrap(err, errProviderNotRetrieved)
	}

	if provider.GetCredentialsSecretReference() == nil {
		return nil, errors.New(errProviderSecretNil)
	}

	secret := &v1.Secret{}
	n := types.NamespacedName{Namespace: provider.Spec.CredentialsSecretRef.Namespace, Name: provider.Spec.CredentialsSecretRef.Name}
	if err := c.kube.Get(ctx, n, secret); err != nil {
		return nil, errors.Wrap(err, errProviderSecretNotRetrieved)
	}

	s, err := c.newServiceFn()
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &lkeClusterPoolExternal{kube: c.kube, service: s}, nil
}

type lkeClusterPoolExternal struct {
	kube    client.Client
	service interface{}
}

func (c *lkeClusterPoolExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.LKEClusterPool)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotLKEClusterPool)
	}

	fmt.Printf("Observing: %+v", cr)

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
		// ConnectionDetails: getConnectionDetails(cr, instance),
	}, nil
}

func (c *lkeClusterPoolExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.LKEClusterPool)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotLKEClusterPool)
	}

	fmt.Printf("Creating: %+v", cr)

	return managed.ExternalCreation{}, nil
}

func (c *lkeClusterPoolExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.LKEClusterPool)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotLKEClusterPool)
	}

	fmt.Printf("Updating: %+v", cr)

	return managed.ExternalUpdate{}, nil
}

func (c *lkeClusterPoolExternal) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.LKEClusterPool)
	if !ok {
		return errors.New(errNotLKEClusterPool)
	}

	fmt.Printf("Deleting: %+v", cr)

	return nil
}
