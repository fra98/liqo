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

package ipctrl

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

// IPReconciler reconciles a Ip object.
type IPReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ipam.liqo.io,resources=ips,verbs=get;list;watch
// +kubebuilder:rbac:groups=ipam.liqo.io,resources=ips/status,verbs=get;list;watch;create;update;patch;delete

// Reconcile Ip objects.
func (r *IPReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Infof("reconciling Ip %q", req.NamespacedName)
	// Fetch the Ip instance
	var nw ipamv1alpha1.IP
	if err := r.Get(ctx, req.NamespacedName, &nw); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(4).Infof("ip %q not found", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		klog.Errorf("an error occurred while getting ip %q: %v", req.NamespacedName, err)
		return ctrl.Result{}, err
	}

	// check if label exists
	remoteclusterID, found := nw.Labels[consts.RemoteClusterID]
	if !found {
		err := fmt.Errorf("ip %q has no remote cluster label", req.NamespacedName)
		klog.Error(err)
		return ctrl.Result{}, err
	}

	desiredIP := nw.Spec.IP
	remappedIP := fakeRemappedIP(nw.Spec.IP, remoteclusterID)

	// Update status
	nw.Status.IP = remappedIP
	if err := r.Client.Status().Update(ctx, &nw); err != nil {
		klog.Errorf("error while updating Ip %q status: %v", req.NamespacedName, err)
		return ctrl.Result{}, err
	}
	klog.Infof("updated Ip %q (%s -> %s)", req.NamespacedName, desiredIP, remappedIP)

	return ctrl.Result{}, nil
}

// SetupWithManager monitors Ip resources.
func (r *IPReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ipamv1alpha1.IP{}).
		Complete(r)
}

func fakeRemappedIP(desiredIP, remoteclusterID string) string {
	_ = remoteclusterID
	return fmt.Sprintf("%s_FAKE", desiredIP)
}
