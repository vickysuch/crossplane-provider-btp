//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/crossplane-contrib/xp-testing/pkg/resources"
	meta "github.com/sap/crossplane-provider-btp/apis"
	"github.com/sap/crossplane-provider-btp/apis/environment/v1alpha1"
	res "sigs.k8s.io/e2e-framework/klient/k8s/resources"

	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestCloudFoundryEnvironment(t *testing.T) {
	var manifestDir = "testdata/crs/cloudfoundry_env"
	var cfName = "cloudfoundry-environment"

	crudFeature := features.New("BTP CF Environment Controller").
		Setup(
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				resources.ImportResources(ctx, t, cfg, manifestDir)
				r, _ := res.New(cfg.Client().RESTConfig())
				_ = meta.AddToScheme(r.GetScheme())
				return ctx
			},
		).
		Assess(
			"Await resources to become synced",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				resources.WaitForResourcesToBeSynced(ctx, cfg, manifestDir, wait.WithTimeout(time.Minute*25))
				return ctx
			},
		).
		Assess(
			"Update with Manager as Technical User",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

				cf := &v1alpha1.CloudFoundryEnvironment{}
				MustGetResource(t, cfg, cfName, nil, cf)

				// managers are initialized with the user that has been used for creation, we should not try to update this again
				newManager := getUserNameFromSecretOrError(t)
				cf.Spec.ForProvider.Managers = []string{newManager}

				resources.AwaitResourceUpdateOrError(ctx, t, cfg, cf)

				return ctx
			},
		).
		Assess(
			"Check Resources Delete",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				resources.DeleteResources(ctx, t, cfg, manifestDir, wait.WithTimeout(time.Minute*25))
				return ctx
			},
		).
		Teardown(resources.DumpManagedResources).
		Feature()

	testenv.Test(t, crudFeature)
}
