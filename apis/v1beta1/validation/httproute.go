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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	gatewayv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

var (
	// repeatableHTTPRouteFilters are filter types that can are allowed to be
	// repeated multiple times in a rule.
	repeatableHTTPRouteFilters = []gatewayv1b1.HTTPRouteFilterType{
		gatewayv1b1.HTTPRouteFilterExtensionRef,
	}

	invalidPathSequences = []string{"//", "/./", "/../", "%2f", "%2F", "#"}
	invalidPathSuffixes  = []string{"/..", "/."}
)

// ValidateHTTPRoute validates HTTPRoute according to the Gateway API specification.
// For additional details of the HTTPRoute spec, refer to:
// https://gateway-api.sigs.k8s.io/v1beta1/references/spec/#gateway.networking.k8s.io/v1beta1.HTTPRoute
func ValidateHTTPRoute(route *gatewayv1b1.HTTPRoute) field.ErrorList {
	return validateHTTPRouteSpec(&route.Spec, field.NewPath("spec"))
}

// ValidateHTTPRouteSpec validates that required fields of spec are set according to the
// HTTPRoute specification.
func ValidateHTTPRouteSpec(spec *gatewayv1b1.HTTPRouteSpec, path *field.Path) field.ErrorList {
	return validateHTTPRouteSpec(spec, path)
}

// validateHTTPRouteSpec validates that required fields of spec are set according to the
// HTTPRoute specification.
func validateHTTPRouteSpec(spec *gatewayv1b1.HTTPRouteSpec, path *field.Path) field.ErrorList {
	return validateHTTPRouteUniqueFilters(spec.Rules, path.Child("rules"))
}

// validateHTTPRouteUniqueFilters validates whether each core and extended filter
// is used at most once in each rule.
func validateHTTPRouteUniqueFilters(rules []gatewayv1b1.HTTPRouteRule, path *field.Path) field.ErrorList {
	var errs field.ErrorList

	for i, rule := range rules {
		counts := map[gatewayv1b1.HTTPRouteFilterType]int{}
		for _, filter := range rule.Filters {
			counts[filter.Type]++
		}
		// custom filters don't have any validation
		for _, key := range repeatableHTTPRouteFilters {
			delete(counts, key)
		}

		for filterType, count := range counts {
			if count > 1 {
				errs = append(errs, field.Invalid(path.Index(i).Child("filters"), filterType, "cannot be used multiple times in the same rule"))
			}
		}

		if errList := validateHTTPBackendUniqueFilters(rule.BackendRefs, path, i); len(errList) > 0 {
			errs = append(errs, errList...)
		}
	}

	return errs
}

func validateHTTPBackendUniqueFilters(ref []gatewayv1b1.HTTPBackendRef, path *field.Path, i int) field.ErrorList {
	var errs field.ErrorList

	for _, bkr := range ref {
		counts := map[gatewayv1b1.HTTPRouteFilterType]int{}
		for _, filter := range bkr.Filters {
			counts[filter.Type]++
		}

		for _, key := range repeatableHTTPRouteFilters {
			delete(counts, key)
		}

		for filterType, count := range counts {
			if count > 1 {
				errs = append(errs, field.Invalid(path.Index(i).Child("BackendRefs"), filterType, "cannot be used multiple times in the same backend"))
			}
		}
	}
	return errs
}

// webhook validation of HTTPPathMatch
func validateHTTPPathMatch(path *gatewayv1b1.HTTPPathMatch, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if path.Type == nil {
		return append(allErrs, field.Required(fldPath.Child("pathType"), "pathType must be specified"))
	}

	if path.Value == nil {
		return append(allErrs, field.Required(fldPath.Child("pathValue"), "pathValue must not be nil."))
	}

	switch *path.Type {
	case gatewayv1b1.PathMatchExact, gatewayv1b1.PathMatchPathPrefix:
		if !strings.HasPrefix(*path.Value, "/") {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("path"), path, "must be an absolute path"))
		}
		if len(*path.Value) > 0 {
			for _, invalidSeq := range invalidPathSequences {
				if strings.Contains(*path.Value, invalidSeq) {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("path"), path, fmt.Sprintf("must not contain '%s'", invalidSeq)))
				}
			}

			for _, invalidSuff := range invalidPathSuffixes {
				if strings.HasSuffix(*path.Value, invalidSuff) {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("path"), path, fmt.Sprintf("cannot end with '%s'", invalidSuff)))
				}
			}
		}
	case gatewayv1b1.PathMatchRegularExpression:
	default:
		pathTypes := []string{string(gatewayv1b1.PathMatchExact), string(gatewayv1b1.PathMatchPathPrefix), string(gatewayv1b1.PathMatchRegularExpression)}
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("pathType"), *path.Type, pathTypes))
	}
	return allErrs
}
