/*
Copyright 2021 The Kubernetes Authors.

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

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/validation/field"

	gatewayv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func TestValidateHTTPRoute(t *testing.T) {
	testService := gatewayv1a2.ObjectName("test-service")

	tests := []struct {
		name     string
		rules    []gatewayv1a2.HTTPRouteRule
		errCount int
	}{{
		name:     "valid httpRoute with no filters",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{
			{
				Matches: []gatewayv1a2.HTTPRouteMatch{
					{
						Path: &gatewayv1a2.HTTPPathMatch{
							Type:  ptrTo(gatewayv1b1.PathMatchType("PathPrefix")),
							Value: ptrTo("/"),
						},
					},
				},
				BackendRefs: []gatewayv1a2.HTTPBackendRef{
					{
						BackendRef: gatewayv1a2.BackendRef{
							BackendObjectReference: gatewayv1a2.BackendObjectReference{
								Name: testService,
								Port: ptrTo(gatewayv1b1.PortNumber(8080)),
							},
							Weight: ptrTo(int32(100)),
						},
					},
				},
			},
		},
	}, {
		name:     "valid httpRoute with 1 filter",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{
			{
				Matches: []gatewayv1a2.HTTPRouteMatch{
					{
						Path: &gatewayv1a2.HTTPPathMatch{
							Type:  ptrTo(gatewayv1b1.PathMatchType("PathPrefix")),
							Value: ptrTo("/"),
						},
					},
				},
				Filters: []gatewayv1a2.HTTPRouteFilter{
					{
						Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
						RequestMirror: &gatewayv1a2.HTTPRequestMirrorFilter{
							BackendRef: gatewayv1a2.BackendObjectReference{
								Name: testService,
								Port: ptrTo(gatewayv1b1.PortNumber(8081)),
							},
						},
					},
				},
			},
		},
	}, {
		name:     "invalid httpRoute with mix of filters and one duplicate",
		errCount: 1,
		rules: []gatewayv1a2.HTTPRouteRule{
			{
				Matches: []gatewayv1a2.HTTPRouteMatch{
					{
						Path: &gatewayv1a2.HTTPPathMatch{
							Type:  ptrTo(gatewayv1b1.PathMatchType("PathPrefix")),
							Value: ptrTo("/"),
						},
					},
				},
				Filters: []gatewayv1a2.HTTPRouteFilter{
					{
						Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
						RequestHeaderModifier: &gatewayv1a2.HTTPRequestHeaderFilter{
							Set: []gatewayv1a2.HTTPHeader{
								{
									Name:  "special-header",
									Value: "foo",
								},
							},
						},
					},
					{
						Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
						RequestMirror: &gatewayv1a2.HTTPRequestMirrorFilter{
							BackendRef: gatewayv1a2.BackendObjectReference{
								Name: testService,
								Port: ptrTo(gatewayv1b1.PortNumber(8080)),
							},
						},
					},
					{
						Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
						RequestHeaderModifier: &gatewayv1a2.HTTPRequestHeaderFilter{
							Add: []gatewayv1a2.HTTPHeader{
								{
									Name:  "my-header",
									Value: "bar",
								},
							},
						},
					},
				},
			},
		},
	}, {
		name:     "valid httpRoute with duplicate ExtensionRef filters",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{
			{
				Matches: []gatewayv1a2.HTTPRouteMatch{
					{
						Path: &gatewayv1a2.HTTPPathMatch{
							Type:  ptrTo(gatewayv1b1.PathMatchType("PathPrefix")),
							Value: ptrTo("/"),
						},
					},
				},
				Filters: []gatewayv1a2.HTTPRouteFilter{
					{
						Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
						RequestHeaderModifier: &gatewayv1a2.HTTPRequestHeaderFilter{
							Set: []gatewayv1a2.HTTPHeader{
								{
									Name:  "special-header",
									Value: "foo",
								},
							},
						},
					},
					{
						Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
						RequestMirror: &gatewayv1a2.HTTPRequestMirrorFilter{
							BackendRef: gatewayv1a2.BackendObjectReference{
								Name: testService,
								Port: ptrTo(gatewayv1b1.PortNumber(8080)),
							},
						},
					},
					{
						Type: "ExtensionRef",
						ExtensionRef: &gatewayv1a2.LocalObjectReference{
							Kind: "Service",
							Name: "test",
						},
					},
					{
						Type: "ExtensionRef",
						ExtensionRef: &gatewayv1a2.LocalObjectReference{
							Kind: "Service",
							Name: "test",
						},
					},
					{
						Type: "ExtensionRef",
						ExtensionRef: &gatewayv1a2.LocalObjectReference{
							Kind: "Service",
							Name: "test",
						},
					},
				},
			},
		},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var errs field.ErrorList
			route := gatewayv1a2.HTTPRoute{Spec: gatewayv1a2.HTTPRouteSpec{Rules: tc.rules}}
			errs = ValidateHTTPRoute(&route)
			if len(errs) != tc.errCount {
				t.Errorf("got %d errors, want %d errors: %s", len(errs), tc.errCount, errs)
			}
		})
	}
}

func TestValidateHTTPBackendUniqueFilters(t *testing.T) {
	var testService gatewayv1a2.ObjectName = "testService"
	tests := []struct {
		name     string
		rules    []gatewayv1a2.HTTPRouteRule
		errCount int
	}{{
		name:     "valid httpRoute Rules backendref filters",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{{
			BackendRefs: []gatewayv1a2.HTTPBackendRef{
				{
					BackendRef: gatewayv1a2.BackendRef{
						BackendObjectReference: gatewayv1a2.BackendObjectReference{
							Name: testService,
							Port: ptrTo(gatewayv1b1.PortNumber(8080)),
						},
						Weight: ptrTo(int32(100)),
					},
					Filters: []gatewayv1a2.HTTPRouteFilter{
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1a2.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1a2.BackendObjectReference{
									Name: testService,
									Port: ptrTo(gatewayv1b1.PortNumber(8080)),
								},
							},
						},
					},
				},
			},
		}},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := gatewayv1a2.HTTPRoute{Spec: gatewayv1a2.HTTPRouteSpec{Rules: tc.rules}}
			errs := ValidateHTTPRoute(&route)
			if len(errs) != tc.errCount {
				t.Errorf("got %d errors, want %d errors: %s", len(errs), tc.errCount, errs)
			}
		})
	}
}

func TestValidateHTTPPathMatch(t *testing.T) {
	tests := []struct {
		name     string
		path     *gatewayv1a2.HTTPPathMatch
		errCount int
	}{{
		name: "invalid httpRoute prefix",
		path: &gatewayv1a2.HTTPPathMatch{
			Type:  ptrTo(gatewayv1b1.PathMatchType("PathPrefix")),
			Value: ptrTo("/"),
		},
		errCount: 0,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := gatewayv1a2.HTTPRoute{Spec: gatewayv1a2.HTTPRouteSpec{
				Rules: []gatewayv1a2.HTTPRouteRule{{
					Matches: []gatewayv1a2.HTTPRouteMatch{{
						Path: tc.path,
					}},
					BackendRefs: []gatewayv1a2.HTTPBackendRef{{
						BackendRef: gatewayv1a2.BackendRef{
							BackendObjectReference: gatewayv1a2.BackendObjectReference{
								Name: gatewayv1a2.ObjectName("test"),
								Port: ptrTo(gatewayv1b1.PortNumber(8080)),
							},
						},
					}},
				}},
			}}

			errs := ValidateHTTPRoute(&route)
			if len(errs) != tc.errCount {
				t.Errorf("got %d errors, want %d errors: %s", len(errs), tc.errCount, errs)
			}
		})
	}
}

func TestValidateHTTPHeaderMatches(t *testing.T) {
	tests := []struct {
		name          string
		headerMatches []gatewayv1a2.HTTPHeaderMatch
		expectErr     string
	}{{
		name:          "no header matches",
		headerMatches: nil,
		expectErr:     "",
	}, {
		name: "no header matched more than once",
		headerMatches: []gatewayv1a2.HTTPHeaderMatch{
			{Name: "Header-Name-1", Value: "val-1"},
			{Name: "Header-Name-2", Value: "val-2"},
			{Name: "Header-Name-3", Value: "val-3"},
		},
		expectErr: "",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := gatewayv1a2.HTTPRoute{Spec: gatewayv1a2.HTTPRouteSpec{
				Rules: []gatewayv1a2.HTTPRouteRule{{
					Matches: []gatewayv1a2.HTTPRouteMatch{{
						Headers: tc.headerMatches,
					}},
					BackendRefs: []gatewayv1a2.HTTPBackendRef{{
						BackendRef: gatewayv1a2.BackendRef{
							BackendObjectReference: gatewayv1a2.BackendObjectReference{
								Name: gatewayv1a2.ObjectName("test"),
								Port: ptrTo(gatewayv1b1.PortNumber(8080)),
							},
						},
					}},
				}},
			}}

			errs := ValidateHTTPRoute(&route)
			if len(tc.expectErr) == 0 {
				assert.Emptyf(t, errs, "expected no errors, got %d errors: %s", len(errs), errs)
			} else {
				require.Lenf(t, errs, 1, "expected one error, got %d errors: %s", len(errs), errs)
				assert.Equal(t, tc.expectErr, errs[0].Error())
			}
		})
	}
}

func TestValidateHTTPQueryParamMatches(t *testing.T) {
	tests := []struct {
		name              string
		queryParamMatches []gatewayv1a2.HTTPQueryParamMatch
		expectErr         string
	}{{
		name:              "no query param matches",
		queryParamMatches: nil,
		expectErr:         "",
	}, {
		name: "no query param matched more than once",
		queryParamMatches: []gatewayv1a2.HTTPQueryParamMatch{
			{Name: "query-param-1", Value: "val-1"},
			{Name: "query-param-2", Value: "val-2"},
			{Name: "query-param-3", Value: "val-3"},
		},
		expectErr: "",
	}, {
		name: "query param names with different casing are not considered duplicates",
		queryParamMatches: []gatewayv1a2.HTTPQueryParamMatch{
			{Name: "query-param-1", Value: "val-1"},
			{Name: "query-param-2", Value: "val-2"},
			{Name: "QUERY-PARAM-1", Value: "val-3"},
		},
		expectErr: "",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := gatewayv1a2.HTTPRoute{Spec: gatewayv1a2.HTTPRouteSpec{
				Rules: []gatewayv1a2.HTTPRouteRule{{
					Matches: []gatewayv1a2.HTTPRouteMatch{{
						QueryParams: tc.queryParamMatches,
					}},
					BackendRefs: []gatewayv1a2.HTTPBackendRef{{
						BackendRef: gatewayv1a2.BackendRef{
							BackendObjectReference: gatewayv1a2.BackendObjectReference{
								Name: gatewayv1a2.ObjectName("test"),
								Port: ptrTo(gatewayv1b1.PortNumber(8080)),
							},
						},
					}},
				}},
			}}

			errs := ValidateHTTPRoute(&route)
			if len(tc.expectErr) == 0 {
				assert.Emptyf(t, errs, "expected no errors, got %d errors: %s", len(errs), errs)
			} else {
				require.Lenf(t, errs, 1, "expected one error, got %d errors: %s", len(errs), errs)
				assert.Equal(t, tc.expectErr, errs[0].Error())
			}
		})
	}
}

func TestValidateServicePort(t *testing.T) {
	portPtr := func(n int) *gatewayv1a2.PortNumber {
		p := gatewayv1a2.PortNumber(n)
		return &p
	}

	groupPtr := func(g string) *gatewayv1a2.Group {
		p := gatewayv1a2.Group(g)
		return &p
	}

	kindPtr := func(k string) *gatewayv1a2.Kind {
		p := gatewayv1a2.Kind(k)
		return &p
	}

	tests := []struct {
		name     string
		rules    []gatewayv1a2.HTTPRouteRule
		errCount int
	}{{
		name:     "default groupkind with port",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{{
			BackendRefs: []gatewayv1a2.HTTPBackendRef{{
				BackendRef: gatewayv1a2.BackendRef{
					BackendObjectReference: gatewayv1a2.BackendObjectReference{
						Name: "backend",
						Port: portPtr(99),
					},
				},
			}},
		}},
	}, {
		name:     "explicit service with port",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{{
			BackendRefs: []gatewayv1a2.HTTPBackendRef{{
				BackendRef: gatewayv1a2.BackendRef{
					BackendObjectReference: gatewayv1a2.BackendObjectReference{
						Group: groupPtr(""),
						Kind:  kindPtr("Service"),
						Name:  "backend",
						Port:  portPtr(99),
					},
				},
			}},
		}},
	}, {
		name:     "explicit ref with no port",
		errCount: 0,
		rules: []gatewayv1a2.HTTPRouteRule{{
			BackendRefs: []gatewayv1a2.HTTPBackendRef{{
				BackendRef: gatewayv1a2.BackendRef{
					BackendObjectReference: gatewayv1a2.BackendObjectReference{
						Group: groupPtr("foo.example.com"),
						Kind:  kindPtr("Foo"),
						Name:  "backend",
					},
				},
			}},
		}},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := gatewayv1a2.HTTPRoute{Spec: gatewayv1a2.HTTPRouteSpec{Rules: tc.rules}}
			errs := ValidateHTTPRoute(&route)
			if len(errs) != tc.errCount {
				t.Errorf("got %d errors, want %d errors: %s", len(errs), tc.errCount, errs)
			}
		})
	}
}

func TestValidateHTTPRouteTypeMatchesField(t *testing.T) {
	tests := []struct {
		name        string
		routeFilter gatewayv1a2.HTTPRouteFilter
		errCount    int
	}{{
		name: "valid HTTPRouteFilterRequestHeaderModifier route filter",
		routeFilter: gatewayv1a2.HTTPRouteFilter{
			Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
			RequestHeaderModifier: &gatewayv1a2.HTTPRequestHeaderFilter{
				Set:    []gatewayv1a2.HTTPHeader{{Name: "name"}},
				Add:    []gatewayv1a2.HTTPHeader{{Name: "add"}},
				Remove: []string{"remove"},
			},
		},
		errCount: 0,
	}, {
		name: "valid HTTPRouteFilterRequestMirror route filter",
		routeFilter: gatewayv1a2.HTTPRouteFilter{
			Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
			RequestMirror: &gatewayv1a2.HTTPRequestMirrorFilter{BackendRef: gatewayv1a2.BackendObjectReference{
				Group:     new(gatewayv1a2.Group),
				Kind:      new(gatewayv1a2.Kind),
				Name:      "name",
				Namespace: new(gatewayv1a2.Namespace),
				Port:      ptrTo(gatewayv1b1.PortNumber(22)),
			}},
		},
		errCount: 0,
	}, {
		name: "valid HTTPRouteFilterExtensionRef filter",
		routeFilter: gatewayv1a2.HTTPRouteFilter{
			Type: gatewayv1b1.HTTPRouteFilterExtensionRef,
			ExtensionRef: &gatewayv1a2.LocalObjectReference{
				Group: "group",
				Kind:  "kind",
				Name:  "name",
			},
		},
		errCount: 0,
	}, {
		name:        "empty type filter is valid (caught by CRD validation)",
		routeFilter: gatewayv1a2.HTTPRouteFilter{},
		errCount:    0,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := gatewayv1a2.HTTPRoute{
				Spec: gatewayv1a2.HTTPRouteSpec{
					Rules: []gatewayv1a2.HTTPRouteRule{{
						Filters: []gatewayv1a2.HTTPRouteFilter{tc.routeFilter},
						BackendRefs: []gatewayv1a2.HTTPBackendRef{{
							BackendRef: gatewayv1a2.BackendRef{
								BackendObjectReference: gatewayv1a2.BackendObjectReference{
									Name: gatewayv1a2.ObjectName("test"),
									Port: ptrTo(gatewayv1b1.PortNumber(8080)),
								},
							},
						}},
					}},
				},
			}
			errs := ValidateHTTPRoute(&route)
			if len(errs) != tc.errCount {
				t.Errorf("got %d errors, want %d errors: %s", len(errs), tc.errCount, errs)
			}
		})
	}
}
