//go:build e2e

package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sap/crossplane-provider-btp/apis"
	"github.com/sap/crossplane-provider-btp/apis/account/v1beta1"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/wait"

	servicemanager "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-service-manager-api-go/pkg"

	"github.com/crossplane-contrib/xp-testing/pkg/resources"
	"github.com/sap/crossplane-provider-btp/apis/account/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	res "sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

var (
	cisCreateName = "e2e-cis-created"

	smImportName  = "e2e-sm-cis-import"
	cisImportName = "e2e-cis-imported"

	smPartialImportName  = "e2e-sm-cis-partial-import"
	cisPartialImportName = "e2e-cis-partial-imported"
)

const (
	cCtxInstanceID = "INSTANCE_ID"
	cCtxBindingID  = "BINDING_ID"
)

func TestCloudManagementCreationFlow(t *testing.T) {
	crudFeatureSuite := features.New("CloudManagement Creation Flow").
		Setup(
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				resources.ImportResources(ctx, t, cfg, "testdata/crs/cloudmanagement/create_flow")
				r, _ := res.New(cfg.Client().RESTConfig())
				_ = apis.AddToScheme(r.GetScheme())

				cm := v1alpha1.CloudManagement{
					ObjectMeta: metav1.ObjectMeta{Name: cisCreateName, Namespace: cfg.Namespace()},
				}
				waitForResource(&cm, cfg, t, wait.WithTimeout(10*time.Minute))
				return ctx
			},
		).
		Assess(
			"Check CloudManagement Resources are fully created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				cm := &v1alpha1.CloudManagement{}
				MustGetResource(t, cfg, cisCreateName, nil, cm)
				// Status bound?
				if cm.Status.AtProvider.Status != v1alpha1.CisStatusBound {
					t.Error("Binding status not set as expected")
				}

				assertProperSecretWritten(t, ctx, cfg, cm)

				// all external resources exist?
				sm := &v1beta1.ServiceManager{}
				MustGetResource(t, cfg, cm.Spec.ForProvider.ServiceManagerRef.Name, nil, sm)

				return ctx
			},
		).Assess(
		"Properly delete all resources", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// k8s resource cleaned up?
			cm := &v1alpha1.CloudManagement{}
			MustGetResource(t, cfg, cisCreateName, nil, cm)
			AwaitResourceDeletionOrFail(ctx, t, cfg, cm, wait.WithTimeout(time.Minute*10))

			return ctx
		},
	).Teardown(
		func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			DeleteResourcesIgnoreMissing(ctx, t, cfg, "cloudmanagement/create_flow", wait.WithTimeout(time.Minute*10))
			return ctx
		},
	).Feature()

	testenv.Test(t, crudFeatureSuite)
}

// Follow Up on that
// couldn't find a way for now to apply a v1beta1.ServiceManager

//func TestCloudManagementImport(t *testing.T) {
//	crudFeatureSuite := features.New("CloudManagement Full Import").
//		Setup(
//			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
//				r, _ := res.New(cfg.Client().RESTConfig())
//				_ = apis.AddToScheme(r.GetScheme())
//
//				importResourcesFromDir(ctx, t, cfg, "cloudmanagement/import/environment")
//
//				// ServiceManager is the requirement for the next steps, so we wait for it to be healthy
//				waitForResource(&v1alpha1.ServiceManager{
//					ObjectMeta: metav1.ObjectMeta{Name: smImportName, Namespace: cfg.Namespace()},
//				}, cfg, t)
//
//				sm := &v1alpha1.ServiceManager{}
//				MustGetResource(t, cfg, smImportName, nil, sm)
//
//				apiClient := configureServiceManagerAPIClient(t, cfg, sm)
//				instanceID := createAPIInstance(t, apiClient, cisImportName)
//
//				if instanceID == nil {
//					t.Errorf("Creating API instance for import testing failed")
//					t.FailNow()
//				}
//				bindingID := createAPIBinding(t, apiClient, cisImportName, instanceID)
//
//				return context.WithValue(
//					context.WithValue(ctx, cCtxInstanceID, *instanceID),
//					cCtxBindingID, *bindingID,
//				)
//			},
//		).
//		Assess(
//			"Check CloudManagement is healthy and binding is imported", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
//				instanceID := ctx.Value(cCtxInstanceID).(string)
//				bindingID := ctx.Value(cCtxBindingID).(string)
//
//				r, _ := res.New(cfg.Client().RESTConfig())
//				_ = apis.AddToScheme(r.GetScheme())
//
//				mutateFn := func(obj k8s.Object) error {
//					crossplane_meta.SetExternalName(obj, instanceID+"/"+bindingID)
//					return nil
//				}
//				createK8sResources(ctx, t, cfg, r, "cloudmanagement/import", "cloudmanagement.yaml", mutateFn)
//
//				waitForResource(&v1alpha1.CloudManagement{
//					ObjectMeta: metav1.ObjectMeta{Name: cisImportName, Namespace: cfg.Namespace()},
//				}, cfg, t)
//
//				cm := &v1alpha1.CloudManagement{}
//				MustGetResource(t, cfg, cisImportName, nil, cm)
//				// Status bound?
//				if cm.Status.AtProvider.Status != v1alpha1.CisStatusBound {
//					t.Error("Binding status not set as expected")
//				}
//				assertProperSecretWritten(t, ctx, cfg, cm)
//				return ctx
//			},
//		).Teardown(
//		func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
//			cis := &v1alpha1.CloudManagement{}
//			MustGetResource(t, cfg, cisImportName, nil, cis)
//			resources.AwaitResourceDeletionOrFail(ctx, t, cfg, cis)
//
//			resources.DeleteResources(ctx, t, cfg, "cloudmanagement/import/environment", wait.WithTimeout(time.Minute*2), false)
//			return ctx
//		},
//	).Feature()
//
//	testenv.Test(t, crudFeatureSuite)
//}
//
//func TestCloudManagementPartialImport(t *testing.T) {
//	crudFeatureSuite := features.New("CloudManagement Partial Import").
//		Setup(
//			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
//				resources.ImportResources(ctx, t, cfg, "cloudmanagement/import_partial/environment")
//				r, _ := res.New(cfg.Client().RESTConfig())
//				_ = apis.AddToScheme(r.GetScheme())
//
//				// ServiceManager is the requirement for the next steps, so we wait for it to be healthy
//				waitForResource(&v1alpha1.ServiceManager{
//					ObjectMeta: metav1.ObjectMeta{Name: smPartialImportName, Namespace: cfg.Namespace()},
//				}, cfg, t)
//
//				sm := &v1alpha1.ServiceManager{}
//				MustGetResource(t, cfg, smPartialImportName, nil, sm)
//
//				apiClient := configureServiceManagerAPIClient(t, cfg, sm)
//				// we just create an instance here, no binding
//				createAPIInstance(t, apiClient, cisPartialImportName)
//				return ctx
//			},
//		).
//		Assess(
//			"Check healthy CloudManagement resource and created binding", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
//				resources.ImportResources(ctx, t, cfg, "cloudmanagement/import_partial/cloudmanagement.yaml")
//				waitForResource(&v1alpha1.CloudManagement{
//					ObjectMeta: metav1.ObjectMeta{Name: cisPartialImportName, Namespace: cfg.Namespace()},
//				}, cfg, t)
//
//				cm := &v1alpha1.CloudManagement{}
//				MustGetResource(t, cfg, cisPartialImportName, nil, cm)
//				// Status bound?
//				if cm.Status.AtProvider.Status != v1alpha1.CisStatusBound {
//					t.Error("Binding status not set as expected")
//				}
//
//				// binding secret properly exposed?
//				assertProperSecretWritten(t, ctx, cfg, cm)
//
//				sm := &v1alpha1.ServiceManager{}
//				MustGetResource(t, cfg, smPartialImportName, nil, sm)
//
//				// binding has been created during import?
//				apiClient := configureServiceManagerAPIClient(t, cfg, sm)
//				assertAPIBindingExists(t, apiClient, cm, true)
//				return ctx
//			},
//		).Teardown(
//		func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
//			cis := &v1alpha1.CloudManagement{}
//			MustGetResource(t, cfg, cisPartialImportName, nil, cis)
//			resources.AwaitResourceDeletionOrFail(ctx, t, cfg, cis)
//
//			resources.DeleteResources(ctx, t, cfg, "cloudmanagement/import_partial/environment", wait.WithTimeout(time.Minute*2), false)
//			return ctx
//		},
//	).Feature()
//
//	testenv.Test(t, crudFeatureSuite)
//}

func importResourcesFromDir(ctx context.Context, t *testing.T, cfg *envconf.Config, crsPath string) {
	fullPath := filepath.Join("./testdata/crs", crsPath)
	if files, err := filepath.Glob(filepath.Join(fullPath, "*.yaml")); err != nil || len(files) < 1 {
		t.Errorf("error while importing resources from %s", fullPath)
	}

	r, _ := resources.GetResourcesWithRESTConfig(cfg)
	_ = apis.AddToScheme(r.GetScheme())
	r.WithNamespace(cfg.Namespace())

	// managed resources fare cluster scoped, so if we patched them with the test namespace it won't do anything
	errdecode := decoder.DecodeEachFile(
		ctx, os.DirFS(fullPath), "*",
		decoder.CreateIgnoreAlreadyExists(r),
	)
	if errdecode != nil {
		t.Fatal(errdecode)
	}
}

func assertProperSecretWritten(t *testing.T, ctx context.Context, cfg *envconf.Config, cm *v1alpha1.CloudManagement) {
	// binding secret written?
	secretName := cm.GetWriteConnectionSecretToReference().Name
	secretNS := cm.GetWriteConnectionSecretToReference().Namespace
	secret := &corev1.Secret{}
	err := cfg.Client().Resources().Get(ctx, secretName, secretNS, secret)
	if err != nil {
		t.Error("Error while loading expected secret from Ref")
	}
	// secret contains correct structure
	if _, ok := secret.Data["uaa.url"]; !ok {
		t.Error("Secret not in proper format")
	}
}

func createAPIInstance(t *testing.T, apiClient *servicemanager.APIClient, externalName string) *string {
	request := apiClient.ServiceInstancesAPI.CreateServiceInstance(context.TODO())
	parameters := map[string]string{"grantType": "clientCredentials"}

	createCisLocalInstanceRequest := servicemanager.CreateServiceInstanceRequestPayload{
		CreateByOfferingAndPlanName: &servicemanager.CreateByOfferingAndPlanName{
			Name:                externalName,
			ServiceOfferingName: "cis",
			ServicePlanName:     "local",
			Parameters:          &parameters,
		},
		CreateByPlanID: nil,
	}

	request = request.CreateServiceInstanceRequestPayload(createCisLocalInstanceRequest)
	request = request.Async(false)
	response, _, err := request.Execute()
	if err != nil {
		t.Errorf("Cannot create cis instance over API")
		return nil
	}
	return response.Id
}

func createAPIBinding(t *testing.T, apiClient *servicemanager.APIClient, externalName string, serviceInstanceId *string) *string {
	request := apiClient.ServiceBindingsAPI.CreateServiceBinding(context.TODO())
	createCisLocalBindingRequest := servicemanager.CreateServiceBindingRequestPayload{
		Name:              externalName,
		ServiceInstanceId: *serviceInstanceId,
		Parameters:        nil,
		BindResource:      nil,
	}
	request = request.CreateServiceBindingRequestPayload(createCisLocalBindingRequest)
	request = request.Async(false)
	res, _, err := request.Execute()

	if err != nil {
		t.Errorf("Cannot create cis binding over API")
	}
	return res.Id
}
