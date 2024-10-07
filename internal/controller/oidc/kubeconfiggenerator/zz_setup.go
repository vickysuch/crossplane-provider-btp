package kubeconfiggenerator

import (
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	oidcapisv1alpha1 "github.com/sap/crossplane-provider-btp/apis/oidc/v1alpha1"
	providerv1alpha1 "github.com/sap/crossplane-provider-btp/apis/v1alpha1"
	"github.com/sap/crossplane-provider-btp/btp"
	"github.com/sap/crossplane-provider-btp/internal/controller/providerconfig"
	"github.com/sap/crossplane-provider-btp/internal/tracking"
)

// Setup adds a controller that reconciles GlobalAccount managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	return providerconfig.DefaultSetup(mgr, o, &oidcapisv1alpha1.KubeConfigGenerator{}, oidcapisv1alpha1.KubeConfigGeneratorGroupKind, oidcapisv1alpha1.KubeConfigGeneratorGroupVersionKind, func(kube client.Client, usage resource.Tracker, resourcetracker tracking.ReferenceResolverTracker, btpNewServiceFn func(cisSecretData []byte, serviceAccountSecretData []byte) (*btp.Client, error)) managed.ExternalConnecter {
		return &connector{
			kube: mgr.GetClient(),
			usage: resource.NewProviderConfigUsageTracker(
				mgr.GetClient(),
				&providerv1alpha1.ProviderConfigUsage{},
			),
			newServiceFn: newKubeConfigClientFn,
		}
	})
}
