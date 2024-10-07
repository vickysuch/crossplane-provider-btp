package example_usage

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	openapi "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-xsuaa-service-api-go/pkg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	GroupComparisionOperator = "equals"
	GroupAttributeName       = "Groups"
)

/**
 * Verifies the whole roundrip CRUD of managed a rolecollection for a user, its the contract we require from the API
 */
func TestAssignUserFlow(t *testing.T) {
	t.Skip("Skipping tests as they are only meant for local testing")

	creds := loadCredentials(t)
	ctx := context.Background()
	client := testClient(creds, ctx)

	origin := "sap.default"
	username := "<EMAIL>"
	roleCollection := "Subaccount Administrator"

	t.Run("User not yet in system", func(t *testing.T) {
		_, res, err := client.UsercontrollerAPI.GetUserByName(ctx, username, origin).Execute()
		expectStatus(t, res, http.StatusNotFound)
		expectError(t, err, true)
	})
	// Run whole workflow only in case user does not already exist in API
	if t.Failed() {
		t.Skip("Skipping remaining test due to unexpected environment")
	}

	t.Run("Assign role to user", func(t *testing.T) {
		_, res, err := client.UsercontrollerAPI.AddRoleCollection(ctx, origin, username, roleCollection).CreateUserIfMissing(true).Execute()
		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Reassigning role shouldn't harm", func(t *testing.T) {
		_, res, err := client.UsercontrollerAPI.AddRoleCollection(ctx, origin, username, roleCollection).CreateUserIfMissing(true).Execute()
		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("User exists with role", func(t *testing.T) {
		data, res, err := client.UsercontrollerAPI.GetUserByName(ctx, username, origin).Execute()
		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
		assert.NotNil(t, data)
		assert.Contains(t, data.RoleCollections, roleCollection)
	})
	t.Run("Revoke Role", func(t *testing.T) {
		_, res, err := client.UsercontrollerAPI.RemoveRoleCollection(ctx, origin, username, roleCollection).Execute()
		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Another Revoke Role returns 404 error", func(t *testing.T) {
		_, res, err := client.UsercontrollerAPI.RemoveRoleCollection(ctx, origin, username, roleCollection).Execute()
		expectStatus(t, res, http.StatusNotFound)
		expectError(t, err, true)
	})
	t.Run("User exists without Role", func(t *testing.T) {
		data, res, err := client.UsercontrollerAPI.GetUserByName(ctx, username, origin).Execute()
		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
		assert.NotNil(t, data)
		assert.NotContains(t, data.RoleCollections, roleCollection)
	})

	// To make test and flow deterministic we remove the user via API again, this is not part of the flow contract
	t.Cleanup(func() {
		cleanUpUser(t, client, ctx, username, origin)
	})
}

/**
 * Verifies the whole roundrip CRUD of managed a rolecollection for a usergroup, its the contract we require from the API
 */
func TestAssignGroupFlow(t *testing.T) {
	t.Skip("Skipping tests as they are only meant for local testing")

	creds := loadCredentials(t)
	ctx := context.Background()
	client := testClient(creds, ctx)

	origin := "aedxvzfuh-platform"
	group := "testgroup"
	roleCollection := "Subaccount Administrator"

	t.Run("No Group yet", func(t *testing.T) {
		res, h, err := client.IdpRoleCollectionAPI.GetIdpAttributeValuesFromRoleCollectionByAttribute(ctx, origin, "Groups", roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		expectGroupMapping(t, res, group, false)
	})
	// Run whole workflow only in case group mapping isn't already there
	if t.Failed() {
		t.Skip("Skipping remaining test due to unexpected environment")
	}
	t.Run("Assign role to group", func(t *testing.T) {
		_, res, err := client.IdpRoleCollectionAPI.AddIdpAttributeToRoleCollection(ctx, origin).IdentityProviderMapping(openapi.IdentityProviderMapping{
			RoleCollectionName: &roleCollection,
			AttributeName:      _ptr(GroupAttributeName),
			AttributeValue:     &group,
			Operator:           _ptr(GroupComparisionOperator),
		}).Execute()
		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Group exists with role", func(t *testing.T) {
		res, h, err := client.IdpRoleCollectionAPI.GetIdpAttributeValuesFromRoleCollection(ctx, origin, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		expectGroupMapping(t, res, group, true)
	})
	t.Run("Revoke Role", func(t *testing.T) {
		_, h, err := client.IdpRoleCollectionAPI.DeleteIdpAttributeToRoleCollection(ctx, origin, GroupAttributeName, GroupComparisionOperator, group, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Another Revoke Role returns 404 error", func(t *testing.T) {
		_, h, err := client.IdpRoleCollectionAPI.DeleteIdpAttributeToRoleCollection(ctx, origin, GroupAttributeName, GroupComparisionOperator, group, roleCollection).Execute()
		expectStatus(t, h, http.StatusNotFound)
		expectError(t, err, true)
	})
	t.Run("No role for group anymore", func(t *testing.T) {
		res, h, err := client.IdpRoleCollectionAPI.GetIdpAttributeValuesFromRoleCollectionByAttribute(ctx, origin, "Groups", roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		expectGroupMapping(t, res, group, false)
	})
}

/**
 * Verifies the whole roundrip CRUD of managed a rolecollection itself, its the contract we require from the API
 */
func TestCRUDFlow(t *testing.T) {
	t.Skip("Skipping tests as they are only meant for local testing")

	creds := loadCredentials(t)
	ctx := context.Background()
	client := testClient(creds, ctx)

	roleCollection := "contract_tests_role"
	viewerRole := openapi.RoleReference{
		RoleTemplateAppId: _ptr("cis-local!b2"),
		RoleTemplateName:  _ptr("Subaccount_Viewer"),
		Name:              _ptr("Subaccount Viewer"),
	}
	adminRole := openapi.RoleReference{
		RoleTemplateAppId: _ptr("cis-local!b2"),
		RoleTemplateName:  _ptr("Subaccount_Admin"),
		Name:              _ptr("Subaccount Admin"),
	}

	t.Run("No Collection yet", func(t *testing.T) {
		_, h, err := client.RolecollectionsAPI.GetRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusNotFound)
		expectError(t, err, true)
	})
	// Run whole workflow only in case collection isn't already there
	if t.Failed() {
		t.Skip("Skipping remaining test due to unexpected environment")
	}
	t.Run("Create roleCollection", func(t *testing.T) {
		_, res, err := client.RolecollectionsAPI.CreateRoleCollection(ctx).RoleCollection(openapi.RoleCollection{
			Name:           roleCollection,
			Description:    _ptr("some description"),
			RoleReferences: []openapi.RoleReference{viewerRole},
		}).Execute()

		expectStatus(t, res, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("RoleCollection exists", func(t *testing.T) {
		res, h, err := client.RolecollectionsAPI.GetRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		expectRoleReference(t, res.RoleReferences, viewerRole, true)
	})
	t.Run("Change description", func(t *testing.T) {
		_, h, err := client.RolecollectionsAPI.ChangeRoleCollectionDescription(ctx, roleCollection).RoleCollectionDescription(openapi.RoleCollectionDescription{Description: _ptr("changed description")}).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Verify changed description", func(t *testing.T) {
		res, h, err := client.RolecollectionsAPI.GetRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		assert.Equal(t, "changed description", *res.Description)
	})
	t.Run("Add Role", func(t *testing.T) {
		_, h, err := client.RolecollectionsAPI.AddRolesToRoleCollection(ctx, roleCollection).RoleReference([]openapi.RoleReference{adminRole}).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Verify added Role", func(t *testing.T) {
		res, h, err := client.RolecollectionsAPI.GetRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		expectRoleReference(t, res.RoleReferences, adminRole, true)
	})
	t.Run("Remove Role", func(t *testing.T) {
		_, h, err := client.RolecollectionsAPI.DeleteRolesFromRoleCollection(ctx, roleCollection).RoleReference([]openapi.RoleReference{viewerRole, adminRole}).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Verify removed Roles", func(t *testing.T) {
		res, h, err := client.RolecollectionsAPI.GetRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
		expectRoleReference(t, res.RoleReferences, viewerRole, false)
		expectRoleReference(t, res.RoleReferences, adminRole, false)
	})
	t.Run("Delete RoleCollection", func(t *testing.T) {
		_, h, err := client.RolecollectionsAPI.DeleteRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusOK)
		expectError(t, err, false)
	})
	t.Run("Verify RoleCollection removed", func(t *testing.T) {
		_, h, err := client.RolecollectionsAPI.GetRoleCollectionByName(ctx, roleCollection).Execute()
		expectStatus(t, h, http.StatusNotFound)
		expectError(t, err, true)
	})
}

func cleanUpUser(t *testing.T, client *openapi.APIClient, ctx context.Context, username, origin string) {
	data, _, err := client.UsercontrollerAPI.GetUserByName(ctx, username, origin).Execute()
	if err != nil {
		t.Errorf("Error while cleaning up user %s: %s", username, err.Error())
		return
	}
	_, _, err = client.UsercontrollerAPI.DeleteUserById(ctx, *data.Id).Execute()
	if err != nil {
		t.Errorf("Error while cleaning up user %s: %s", username, err.Error())
		return
	}
	t.Log("Successfully cleaned up user", username)
}

func loadCredentials(t *testing.T) Env {
	env := Env{
		ClientId:     os.Getenv("XSUAA_CLIENT_ID"),
		ClientSecret: os.Getenv("XSUAA_CLIENT_SECRET"),
		TokenURL:     os.Getenv("XSUAA_TOKEN_URL"),
		Url:          os.Getenv("XSUAA_URL"),
	}
	if !env.fullyConfigured() {
		t.Fatal("Environment not properly configured, failing test")
	}
	return env
}

func testClient(creds Env, ctx context.Context) *openapi.APIClient {
	config := clientcredentials.Config{
		ClientID:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		TokenURL:     creds.TokenURL,
	}

	smURL, _ := url.Parse(creds.Url)

	apiClientConfig := openapi.NewConfiguration()
	apiClientConfig.Host = smURL.Host
	apiClientConfig.Scheme = smURL.Scheme
	apiClientConfig.HTTPClient = config.Client(ctx)

	return openapi.NewAPIClient(apiClientConfig)
}

func expectGroupMapping(t *testing.T, attrs []openapi.RoleCollectionAttribute, group string, expect bool) {
	found := findGroupMapping(attrs, group)
	if expect && found == nil {
		t.Errorf("Expected group %s not found", group)
	}
	if !expect && found != nil {
		t.Errorf("Unexpected group %s found", group)
	}
}

func expectRoleReference(t *testing.T, roleRefs []openapi.RoleReference, wantRoleRef openapi.RoleReference, expect bool) {
	found := findRoleReference(roleRefs, wantRoleRef)
	if expect && found == nil {
		t.Errorf("Expected roleref %s not found", *wantRoleRef.Name)
	}
	if !expect && found != nil {
		t.Errorf("Unexpected roleref %s found", *wantRoleRef.Name)
	}
}

func expectStatus(t *testing.T, res *http.Response, expected int) {
	if res.StatusCode != expected {
		t.Errorf("Enexpected status code %d", res.StatusCode)
	}
}

func expectError(t *testing.T, apiErr error, expectErr bool) {
	if expectErr == (apiErr == nil) {
		t.Errorf("API returned unexpected Error %s", apiErr.Error())
	}
}

func findGroupMapping(attrs []openapi.RoleCollectionAttribute, group string) *openapi.RoleCollectionAttribute {
	if attrs == nil {
		return nil
	}
	for _, a := range attrs {
		if *a.AttributeName == GroupAttributeName && *a.AttributeValue == group && *a.ComparisonOperator == GroupComparisionOperator {
			return &a
		}
	}
	return nil
}

func findRoleReference(roleRefs []openapi.RoleReference, wantRoleRef openapi.RoleReference) *openapi.RoleReference {
	if roleRefs == nil {
		return nil
	}
	for _, rr := range roleRefs {
		if *rr.RoleTemplateAppId == *wantRoleRef.RoleTemplateAppId && *rr.RoleTemplateName == *wantRoleRef.RoleTemplateName && *rr.Name == *wantRoleRef.Name {
			return &rr
		}
	}
	return nil
}

func _ptr(v string) *string {
	return &v
}

type Env struct {
	ClientId     string
	ClientSecret string
	TokenURL     string
	Url          string
}

func (e Env) fullyConfigured() bool {
	return e.ClientId != "" && e.ClientSecret != "" && e.TokenURL != "" && e.Url != ""
}
