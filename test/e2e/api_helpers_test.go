//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/sap/crossplane-provider-btp/apis/account/v1alpha1"
	"github.com/sap/crossplane-provider-btp/apis/account/v1beta1"
	providerv1alpha1 "github.com/sap/crossplane-provider-btp/apis/v1alpha1"
	"github.com/sap/crossplane-provider-btp/btp"
	saas_client "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-saas-provisioning-api-go/pkg"
	servicemanager "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-service-manager-api-go/pkg"
	"golang.org/x/oauth2/clientcredentials"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

// configureSaasProvisioningAPIClient creates are APIClient from a servicemanager instance, meant only for verifying API Objects withing e2e tests
func configureServiceManagerAPIClient(t *testing.T, cfg *envconf.Config, smCr *v1beta1.ServiceManager) *servicemanager.APIClient {
	secretRef := smCr.GetWriteConnectionSecretToReference()

	secret := &v1.Secret{}
	MustGetResource(t, cfg, secretRef.Name, &secretRef.Namespace, secret)

	clientId := string(secret.Data[v1alpha1.ResourceCredentialsClientId])
	clientSecret := string(secret.Data[v1alpha1.ResourceCredentialsClientSecret])
	tokenUrl := string(secret.Data[v1alpha1.ResourceCredentialsXsuaaUrl])
	smUrl := string(secret.Data[v1alpha1.ResourceCredentialsServiceManagerUrl])

	endPointParams := url.Values{}
	endPointParams.Add("grant_type", "client_credentials")

	config := clientcredentials.Config{
		ClientID:       clientId,
		ClientSecret:   clientSecret,
		TokenURL:       tokenUrl + "/oauth/token",
		EndpointParams: endPointParams,
	}

	smURL, err := url.Parse(smUrl)
	if err != nil {
		t.Error("Can not create Servicemanager Client due to misconfiguration")
	}

	apiClientConfig := servicemanager.NewConfiguration()
	apiClientConfig.Host = smURL.Host
	apiClientConfig.Scheme = smURL.Scheme
	apiClientConfig.HTTPClient = config.Client(context.TODO())

	client := servicemanager.NewAPIClient(apiClientConfig)
	return client
}

// configureSaasProvisioningAPIClient creates are APIClient from a cloudmangement instance, meant only for verifying API Objects withing e2e tests
func configureSaasProvisioningAPIClient(t *testing.T, cfg *envconf.Config, sub *v1alpha1.CloudManagement) *saas_client.APIClient {
	secretRef := sub.GetWriteConnectionSecretToReference()

	secret := &v1.Secret{}
	MustGetResource(t, cfg, secretRef.Name, &secretRef.Namespace, secret)

	cisBinding := secret.Data[providerv1alpha1.RawBindingKey]

	var cisCredential btp.CISCredential
	if err := json.Unmarshal(cisBinding, &cisCredential); err != nil {
		t.Error("Error while unwrapping api credentials from created cis instance")
	}

	c := saas_client.NewConfiguration()

	config := &clientcredentials.Config{
		ClientID:     cisCredential.Uaa.Clientid,
		ClientSecret: cisCredential.Uaa.Clientsecret,
		TokenURL:     cisCredential.Uaa.Url + "/oauth/token",
	}
	c.HTTPClient = config.Client(context.TODO())
	c.Servers = []saas_client.ServerConfiguration{{URL: cisCredential.Endpoints.SaasRegistryServiceUrl}}

	return saas_client.NewAPIClient(c)
}
