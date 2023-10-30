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

	"k8s.io/apimachinery/pkg/util/validation/field"
	utilpointer "k8s.io/utils/pointer"

	"sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	pkgutils "sigs.k8s.io/gateway-api/pkg/util"
)

func TestValidateHTTPRoute(t *testing.T) {
	testService := gatewayv1b1.ObjectName("test-service")
	specialService := gatewayv1b1.ObjectName("special-service")
	tests := []struct {
		name     string
		rules    []gatewayv1b1.HTTPRouteRule
		errCount int
	}{
		{
			name: "valid httpRoute with no filters",
			rules: []gatewayv1b1.HTTPRouteRule{
				{
					Matches: []gatewayv1b1.HTTPRouteMatch{
						{
							Path: &gatewayv1b1.HTTPPathMatch{
								Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
								Value: utilpointer.String("/"),
							},
						},
					},
					BackendRefs: []gatewayv1b1.HTTPBackendRef{
						{
							BackendRef: gatewayv1b1.BackendRef{
								BackendObjectReference: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8080),
								},
								Weight: utilpointer.Int32(100),
							},
						},
					},
				},
			},
			errCount: 0,
		},
		{
			name: "valid httpRoute with 1 filter",
			rules: []gatewayv1b1.HTTPRouteRule{
				{
					Matches: []gatewayv1b1.HTTPRouteMatch{
						{
							Path: &gatewayv1b1.HTTPPathMatch{
								Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
								Value: utilpointer.String("/"),
							},
						},
					},
					Filters: []gatewayv1b1.HTTPRouteFilter{
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8081),
								},
							},
						},
					},
				},
			},
			errCount: 0,
		},
		{
			name: "invalid httpRoute with 2 extended filters",
			rules: []gatewayv1b1.HTTPRouteRule{
				{
					Matches: []gatewayv1b1.HTTPRouteMatch{
						{
							Path: &gatewayv1b1.HTTPPathMatch{
								Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
								Value: utilpointer.String("/"),
							},
						},
					},
					Filters: []gatewayv1b1.HTTPRouteFilter{
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: specialService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
					},
				},
			},
			errCount: 1,
		},
		{
			name: "invalid httpRoute with mix of filters and one duplicate",
			rules: []gatewayv1b1.HTTPRouteRule{
				{
					Matches: []gatewayv1b1.HTTPRouteMatch{
						{
							Path: &gatewayv1b1.HTTPPathMatch{
								Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
								Value: utilpointer.String("/"),
							},
						},
					},
					Filters: []gatewayv1b1.HTTPRouteFilter{
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: &gatewayv1b1.HTTPRequestHeaderFilter{
								Set: []gatewayv1b1.HTTPHeader{
									{
										Name:  "special-header",
										Value: "foo",
									},
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: &gatewayv1b1.HTTPRequestHeaderFilter{
								Add: []gatewayv1b1.HTTPHeader{
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
			errCount: 1,
		},
		{
			name: "invalid httpRoute with multiple duplicate filters",
			rules: []gatewayv1b1.HTTPRouteRule{
				{
					Matches: []gatewayv1b1.HTTPRouteMatch{
						{
							Path: &gatewayv1b1.HTTPPathMatch{
								Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
								Value: utilpointer.String("/"),
							},
						},
					},
					Filters: []gatewayv1b1.HTTPRouteFilter{
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: &gatewayv1b1.HTTPRequestHeaderFilter{
								Set: []gatewayv1b1.HTTPHeader{
									{
										Name:  "special-header",
										Value: "foo",
									},
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: &gatewayv1b1.HTTPRequestHeaderFilter{
								Add: []gatewayv1b1.HTTPHeader{
									{
										Name:  "my-header",
										Value: "bar",
									},
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: specialService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
					},
				},
			},
			errCount: 2,
		},
		{
			name: "valid httpRoute with duplicate ExtensionRef filters",
			rules: []gatewayv1b1.HTTPRouteRule{
				{
					Matches: []gatewayv1b1.HTTPRouteMatch{
						{
							Path: &gatewayv1b1.HTTPPathMatch{
								Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
								Value: utilpointer.String("/"),
							},
						},
					},
					Filters: []gatewayv1b1.HTTPRouteFilter{
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: &gatewayv1b1.HTTPRequestHeaderFilter{
								Set: []gatewayv1b1.HTTPHeader{
									{
										Name:  "special-header",
										Value: "foo",
									},
								},
							},
						},
						{
							Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
							RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
								BackendRef: gatewayv1b1.BackendObjectReference{
									Name: testService,
									Port: pkgutils.PortNumberPtr(8080),
								},
							},
						},
						{
							Type: "ExtensionRef",
						},
						{
							Type: "ExtensionRef",
						},
						{
							Type: "ExtensionRef",
						},
					},
				},
			},
			errCount: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateHTTPRouteUniqueFilters(tc.rules, field.NewPath("spec").Child("rules"))
			if len(errs) != tc.errCount {
				t.Errorf("ValidateHTTPRoute() got %v errors, want %v errors", len(errs), tc.errCount)
			}
		})
	}
}

func TestValidateHTTPBackendUniqueFilters(t *testing.T) {
	var testService v1beta1.ObjectName = "testService"
	var specialService v1beta1.ObjectName = "specialService"
	tests := []struct {
		name     string
		hRoute   gatewayv1b1.HTTPRoute
		errCount int
	}{
		{
			name: "valid httpRoute Rules backendref filters",
			hRoute: gatewayv1b1.HTTPRoute{
				Spec: gatewayv1b1.HTTPRouteSpec{
					Rules: []gatewayv1b1.HTTPRouteRule{
						{
							BackendRefs: []gatewayv1b1.HTTPBackendRef{
								{
									BackendRef: gatewayv1b1.BackendRef{
										BackendObjectReference: gatewayv1b1.BackendObjectReference{
											Name: testService,
											Port: pkgutils.PortNumberPtr(8080),
										},
										Weight: utilpointer.Int32(100),
									},
									Filters: []gatewayv1b1.HTTPRouteFilter{
										{
											Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
											RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
												BackendRef: gatewayv1b1.BackendObjectReference{
													Name: testService,
													Port: pkgutils.PortNumberPtr(8080),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			errCount: 0,
		},
		{
			name: "invalid httpRoute Rules backendref filters",
			hRoute: gatewayv1b1.HTTPRoute{
				Spec: gatewayv1b1.HTTPRouteSpec{
					Rules: []gatewayv1b1.HTTPRouteRule{
						{
							BackendRefs: []gatewayv1b1.HTTPBackendRef{
								{
									Filters: []gatewayv1b1.HTTPRouteFilter{
										{
											Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
											RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
												BackendRef: gatewayv1b1.BackendObjectReference{
													Name: testService,
													Port: pkgutils.PortNumberPtr(8080),
												},
											},
										},
										{
											Type: gatewayv1b1.HTTPRouteFilterRequestMirror,
											RequestMirror: &gatewayv1b1.HTTPRequestMirrorFilter{
												BackendRef: gatewayv1b1.BackendObjectReference{
													Name: specialService,
													Port: pkgutils.PortNumberPtr(8080),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			errCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for index, rule := range tc.hRoute.Spec.Rules {
				errs := validateHTTPBackendUniqueFilters(rule.BackendRefs, field.NewPath("spec").Child("rules"), index)
				if len(errs) != tc.errCount {
					t.Errorf("ValidateHTTPRoute() got %d errors, want %d errors", len(errs), tc.errCount)
				}
			}
		})
	}
}

func TestValidateHTTPPathMatch(t *testing.T) {
	tests := []struct {
		name     string
		path     *gatewayv1b1.HTTPPathMatch
		errCount int
	}{
		{
			name: "invalid httpRoute prefix",
			path: &gatewayv1b1.HTTPPathMatch{
				Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
				Value: utilpointer.String("/."),
			},
			errCount: 1,
		},
		{
			name: "invalid httpRoute Exact",
			path: &gatewayv1b1.HTTPPathMatch{
				Type:  pkgutils.PathMatchTypePtr("Exact"),
				Value: utilpointer.String("/foo/./bar"),
			},
			errCount: 1,
		},
		{
			name: "invalid httpRoute prefix",
			path: &gatewayv1b1.HTTPPathMatch{
				Type:  pkgutils.PathMatchTypePtr("PathPrefix"),
				Value: utilpointer.String("/"),
			},
			errCount: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateHTTPPathMatch(tc.path, field.NewPath("spec").Child("rules").Child("matches").Child("path"))
			if len(errs) != tc.errCount {
				t.Errorf("TestValidateHTTPPathMatch() got %v errors, want %v errors", len(errs), tc.errCount)
			}
		})
	}
}
