package example_usage

import (
	"context"
	"fmt"
	"testing"

	openapi "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-saas-provisioning-api-go/pkg"
	"golang.org/x/oauth2/clientcredentials"
)

func Test_openapi_ApplicationsApiService(t *testing.T) {

	configuration := openapi.NewConfiguration()

	config := clientcredentials.Config{
		// credentials from local cis binding (referenced by subaccount in providers)
		ClientID:     "...",
		ClientSecret: "...",
		TokenURL:     "...",
	}

	ctx := context.Background()

	configuration.HTTPClient = config.Client(context.Background())
	configuration.Servers = []openapi.ServerConfiguration{{
		// provisioning service url from local cis binding
		URL: "https://saas-manager.cfapps.eu12.hana.ondemand.com",
	}}

	client := openapi.NewAPIClient(configuration)

	t.Run("Test ProvisioningAPI GetApplications", func(t *testing.T) {
		res, raw, err := client.SubscriptionOperationsForAppConsumersAPI.GetEntitledApplication(ctx, "sapappstudio").PlanName("standard-edition").Execute()
		if err != nil {
			return
		}

		fmt.Println(res)
		fmt.Println(raw)
		fmt.Println(err)

	})

	t.Run("Test ProvisioningAPI CreateApplication", func(t *testing.T) {
		plan := "standard-edition"
		raw, err := client.SubscriptionOperationsForAppConsumersAPI.CreateSubscriptionAsync(ctx, "sapappstudio").CreateSubscriptionRequestPayload(openapi.CreateSubscriptionRequestPayload{PlanName: &plan}).Execute()
		if err != nil {
			return
		}

		fmt.Println(raw)
		fmt.Println(err)
	})

	t.Run("Test ProvisioningAPI DeleteApplication", func(t *testing.T) {
		raw, err := client.SubscriptionOperationsForAppConsumersAPI.DeleteSubscriptionAsync(ctx, "sapappstudio").Execute()
		if err != nil {
			return
		}

		fmt.Println(raw)
		fmt.Println(err)
	})

}
