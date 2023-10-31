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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1alpha2 "sigs.k8s.io/gateway-api/pkg/client/clientset/gateway/versioned/typed/apis/v1alpha2"
)

type FakeGatewayV1alpha2 struct {
	*testing.Fake
}

func (c *FakeGatewayV1alpha2) Gateways(namespace string) v1alpha2.GatewayInterface {
	return &FakeGateways{c, namespace}
}

func (c *FakeGatewayV1alpha2) GatewayClasses() v1alpha2.GatewayClassInterface {
	return &FakeGatewayClasses{c}
}

func (c *FakeGatewayV1alpha2) HTTPRoutes(namespace string) v1alpha2.HTTPRouteInterface {
	return &FakeHTTPRoutes{c, namespace}
}

func (c *FakeGatewayV1alpha2) ReferenceGrants(namespace string) v1alpha2.ReferenceGrantInterface {
	return &FakeReferenceGrants{c, namespace}
}

func (c *FakeGatewayV1alpha2) TCPRoutes(namespace string) v1alpha2.TCPRouteInterface {
	return &FakeTCPRoutes{c, namespace}
}

func (c *FakeGatewayV1alpha2) TLSRoutes(namespace string) v1alpha2.TLSRouteInterface {
	return &FakeTLSRoutes{c, namespace}
}

func (c *FakeGatewayV1alpha2) UDPRoutes(namespace string) v1alpha2.UDPRouteInterface {
	return &FakeUDPRoutes{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeGatewayV1alpha2) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
