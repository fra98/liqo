// Copyright 2019-2025 The Liqo Authors
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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"

	offloadingv1beta1 "github.com/liqotech/liqo/apis/offloading/v1beta1"
	versioned "github.com/liqotech/liqo/pkg/client/clientset/versioned"
	internalinterfaces "github.com/liqotech/liqo/pkg/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/liqotech/liqo/pkg/client/listers/offloading/v1beta1"
)

// VirtualNodeInformer provides access to a shared informer and lister for
// VirtualNodes.
type VirtualNodeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.VirtualNodeLister
}

type virtualNodeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewVirtualNodeInformer constructs a new informer for VirtualNode type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewVirtualNodeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredVirtualNodeInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredVirtualNodeInformer constructs a new informer for VirtualNode type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredVirtualNodeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OffloadingV1beta1().VirtualNodes(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OffloadingV1beta1().VirtualNodes(namespace).Watch(context.TODO(), options)
			},
		},
		&offloadingv1beta1.VirtualNode{},
		resyncPeriod,
		indexers,
	)
}

func (f *virtualNodeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredVirtualNodeInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *virtualNodeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&offloadingv1beta1.VirtualNode{}, f.defaultInformer)
}

func (f *virtualNodeInformer) Lister() v1beta1.VirtualNodeLister {
	return v1beta1.NewVirtualNodeLister(f.Informer().GetIndexer())
}
