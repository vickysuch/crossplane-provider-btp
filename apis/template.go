// Package apis contains Kubernetes API for the Template provider.
package apis

import (
	accountv1alpha1 "github.com/sap/crossplane-provider-btp/apis/account/v1alpha1"
	accountv1beta1 "github.com/sap/crossplane-provider-btp/apis/account/v1beta1"
	environmentv1alpha1 "github.com/sap/crossplane-provider-btp/apis/environment/v1alpha1"
	oidcv1alpha1 "github.com/sap/crossplane-provider-btp/apis/oidc/v1alpha1"
	"github.com/sap/crossplane-provider-btp/apis/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(
		AddToSchemes,
		v1alpha1.SchemeBuilder.AddToScheme,
		accountv1alpha1.SchemeBuilder.AddToScheme,
		environmentv1alpha1.SchemeBuilder.AddToScheme,
		oidcv1alpha1.SchemeBuilder.AddToScheme,
		accountv1beta1.SchemeBuilder.AddToScheme,
	)
}
