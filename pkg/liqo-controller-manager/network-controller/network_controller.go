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

package networkctrl

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ipamv1alpha1 "github.com/liqotech/liqo/apis/ipam/v1alpha1"
	"github.com/liqotech/liqo/pkg/consts"
)

// NetworkReconciler reconciles a Network object.
type NetworkReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ipam.liqo.io,resources=networks,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.liqo.io,resources=networks/status,verbs=get;list;watch;create;update;patch;delete

// Reconcile Network objects.
func (r *NetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Infof("reconciling Network %q", req.NamespacedName)
	// Fetch the Network instance
	var nw ipamv1alpha1.Network
	if err := r.Get(ctx, req.NamespacedName, &nw); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(4).Infof("network %q not found", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		klog.Errorf("an error occurred while getting network %q: %v", req.NamespacedName, err)
		return ctrl.Result{}, err
	}

	// check if label exists
	remoteclusterID, found := nw.Labels[consts.RemoteClusterID]
	if !found {
		err := fmt.Errorf("network %q has no remote cluster label", req.NamespacedName)
		klog.Error(err)
		return ctrl.Result{}, err
	}

	desiredCIDR := nw.Spec.CIDR
	remappedCIDR := fakeNetworkCIDR(nw.Spec.CIDR, remoteclusterID)

	// Update status
	nw.Status.CIDR = remappedCIDR
	if err := r.Client.Status().Update(ctx, &nw); err != nil {
		klog.Errorf("error while updating Network %q status: %v", req.NamespacedName, err)
		return ctrl.Result{}, err
	}
	klog.Infof("updated Network %q (%s -> %s)", req.NamespacedName, desiredCIDR, remappedCIDR)

	return ctrl.Result{}, nil
}

// SetupWithManager monitors Network resources.
func (r *NetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ipamv1alpha1.Network{}).
		Complete(r)
}

func fakeNetworkCIDR(desiredCIDR, remoteclusterID string) string {
	_ = remoteclusterID
	return fmt.Sprintf("%s_FAKE", desiredCIDR)
}
