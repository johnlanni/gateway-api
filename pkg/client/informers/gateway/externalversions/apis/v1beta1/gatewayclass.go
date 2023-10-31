/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	apisv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	versioned "sigs.k8s.io/gateway-api/pkg/client/clientset/gateway/versioned"
	internalinterfaces "sigs.k8s.io/gateway-api/pkg/client/informers/gateway/externalversions/internalinterfaces"
	v1beta1 "sigs.k8s.io/gateway-api/pkg/client/listers/gateway/apis/v1beta1"
)

// GatewayClassInformer provides access to a shared informer and lister for
// GatewayClasses.
type GatewayClassInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.GatewayClassLister
}

type gatewayClassInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewGatewayClassInformer constructs a new informer for GatewayClass type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGatewayClassInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGatewayClassInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredGatewayClassInformer constructs a new informer for GatewayClass type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGatewayClassInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GatewayV1beta1().GatewayClasses().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GatewayV1beta1().GatewayClasses().Watch(context.TODO(), options)
			},
		},
		&apisv1beta1.GatewayClass{},
		resyncPeriod,
		indexers,
	)
}

func (f *gatewayClassInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGatewayClassInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gatewayClassInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apisv1beta1.GatewayClass{}, f.defaultInformer)
}

func (f *gatewayClassInformer) Lister() v1beta1.GatewayClassLister {
	return v1beta1.NewGatewayClassLister(f.Informer().GetIndexer())
}
