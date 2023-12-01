// Copyright 2019-2023 The Liqo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package route

import (
	"context"
	"errors"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	networkingv1alpha1 "github.com/liqotech/liqo/apis/networking/v1alpha1"
)

// RouteConfigurationReconciler manage Configuration lifecycle.
//
//nolint:revive // We usually use the name of the reconciled resource in the controller name.
type RouteConfigurationReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	EventsRecorder record.EventRecorder
	// Labels used to filter the reconciled resources.
	Labels map[string]string
}

// NewRouteConfigurationReconciler returns a new RouteConfigurationReconciler.
func NewRouteConfigurationReconciler(cl client.Client, s *runtime.Scheme,
	er record.EventRecorder, labels map[string]string) (*RouteConfigurationReconciler, error) {
	return &RouteConfigurationReconciler{
		Client:         cl,
		Scheme:         s,
		EventsRecorder: er,
		Labels:         labels,
	}, nil
}

// cluster-role
// +kubebuilder:rbac:groups=networking.liqo.io,resources=routeconfigurations,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=networking.liqo.io,resources=routeconfigurations/status,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=networking.liqo.io,resources=routeconfigurations/finalizers,verbs=update

// Reconcile manage RouteConfigurations, applying nftables configuration.
func (r *RouteConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	routeconfiguration := &networkingv1alpha1.RouteConfiguration{}
	if err = r.Get(ctx, req.NamespacedName, routeconfiguration); err != nil {
		if apierrors.IsNotFound(err) {
			klog.Infof("There is no routeconfiguration %s", req.String())
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to get the routeconfiguration %q: %w", req.NamespacedName, err)
	}

	klog.V(4).Infof("Reconciling routeconfiguration %s", req.String())

	defer func() {
		err = r.UpdateStatus(ctx, r.EventsRecorder, routeconfiguration, err)
	}()

	var tableID uint32
	tableID, err = GetTableID(routeconfiguration.Spec.Table.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Manage Finalizers and routeconfiguration deletion.
	deleting := !routeconfiguration.ObjectMeta.DeletionTimestamp.IsZero()
	containsFinalizer := ctrlutil.ContainsFinalizer(routeconfiguration, routeconfigurationControllerFinalizer)
	switch {
	case !deleting && !containsFinalizer:
		if err = r.ensureRouteConfigurationFinalizerPresence(ctx, routeconfiguration); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil

	case deleting && containsFinalizer:
		for i := range routeconfiguration.Spec.Table.Rules {
			if err = EnsureRuleAbsence(&routeconfiguration.Spec.Table.Rules[i], tableID); err != nil {
				return ctrl.Result{}, err
			}
		}

		if err = EnsureTableAbsence(routeconfiguration, tableID); err != nil {
			return ctrl.Result{}, err
		}

		if err = r.ensureRouteConfigurationFinalizerAbsence(ctx, routeconfiguration); err != nil {
			return ctrl.Result{}, err
		}

		klog.V(2).Infof("RouteConfiguration %s deleted", req.String())

		return ctrl.Result{}, nil

	case deleting && !containsFinalizer:
		return ctrl.Result{}, nil
	}

	if err = CleanRules(routeconfiguration.Spec.Table.Rules, tableID); err != nil {
		return ctrl.Result{}, err
	}

	if err = EnforceTablePresence(routeconfiguration, tableID); err != nil {
		return ctrl.Result{}, err
	}

	for i := range routeconfiguration.Spec.Table.Rules {
		if err = EnsureRulePresence(&routeconfiguration.Spec.Table.Rules[i], tableID); err != nil {
			return ctrl.Result{}, err
		}
	}

	klog.V(4).Infof("Applying routeconfiguration %s", req.String())

	return ctrl.Result{}, nil
}

// SetupWithManager register the RouteConfigurationReconciler to the manager.
func (r *RouteConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	filterByLabelsPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{MatchLabels: r.Labels})
	if err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1alpha1.RouteConfiguration{}, builder.WithPredicates(filterByLabelsPredicate)).
		Complete(r)
}

// UpdateStatus updates the status of the given RouteConfiguration.
func (r *RouteConfigurationReconciler) UpdateStatus(ctx context.Context, er record.EventRecorder,
	routeconfiguration *networkingv1alpha1.RouteConfiguration, err error) error {
	// TODO: implement this function.
	er.Eventf(routeconfiguration, "Normal", "RouteConfigurationUpdate", "RouteConfiguration: %s", "TODO")
	if clerr := r.Client.Status().Update(ctx, routeconfiguration); clerr != nil {
		err = errors.Join(err, clerr)
	}
	return err
}