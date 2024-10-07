package example_usage

import (
	"context"
	"fmt"
	"testing"

	openapi "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-accounts-service-api-go/pkg"
	"golang.org/x/oauth2/clientcredentials"
)

func Test_openapi_EnvironmentsApiService(t *testing.T) {

	configuration := openapi.NewConfiguration()

	config := clientcredentials.Config{
		// credentials from local cis binding (referenced by subaccount in providers)
		//ClientID:     "...",
		//ClientSecret: "...",
		//TokenURL:     "...",
	}

	ctx := context.Background()

	configuration.HTTPClient = config.Client(context.Background())
	configuration.Servers = []openapi.ServerConfiguration{{
		// provisioning service url from local cis binding
		URL: "https://accounts-service.cfapps.sap.hana.ondemand.com",
	}}

	client := openapi.NewAPIClient(configuration)

	t.Run("Test AccountsAPI GetSubaccounts", func(t *testing.T) {

		req := client.SubaccountOperationsAPI.GetSubaccounts(ctx)

		execute, h, err := req.Execute()
		if err != nil {
			return
		}

		fmt.Println(execute)
		fmt.Println(h)
		fmt.Println(err)

	})

}
