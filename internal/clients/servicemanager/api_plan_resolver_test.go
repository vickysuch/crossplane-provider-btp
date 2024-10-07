package servicemanager

import (
	"context"
	"net/http"
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/sap/crossplane-provider-btp/internal"
	servicemanager "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-service-manager-api-go/pkg"
)

func TestNewServiceManagerClient(t *testing.T) {
	tests := []struct {
		name    string
		creds   *BindingCredentials
		success bool
	}{
		{
			name: "Invalid SM URL",
			creds: &BindingCredentials{
				Clientid:     internal.Ptr("someClientId"),
				Clientsecret: internal.Ptr("someClientSecret"),
				SmUrl:        internal.Ptr("::noUrl::"),
				Url:          internal.Ptr("https://valid.url"),
			},
			success: false,
		},
		{
			name:    "Success",
			creds:   &BindingCredentials{},
			success: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewServiceManagerClient(context.TODO(), tc.creds)
			if tc.success != (err == nil) {
				t.Errorf("Unexpected error return; Expected error: %v, Returned: %v", !tc.success, err)
			}
			if tc.success != (client != nil) {
				t.Errorf("Unexpected client return; Returned: %v", client)
			}
		})
	}
}

func TestPlanIDByName(t *testing.T) {
	type args struct {
		listOfferingsMockFn func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error)
		listPlansMockFn     func() (*servicemanager.ServicePlanResponseList, *http.Response, error)
	}
	tests := []struct {
		name string
		args args

		wantErr bool
		wantID  string
	}{
		{
			name: "offeringError",
			args: args{
				listOfferingsMockFn: func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error) {
					return nil, nil, errors.New("offeringApiError")
				},
			},
			wantErr: true,
		},
		{
			name: "offeringNotFound",
			args: args{
				listOfferingsMockFn: func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error) {
					return &servicemanager.ServiceOfferingResponseList{
						Items: []servicemanager.ServiceOfferingResponseObject{},
					}, nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "plansError",
			args: args{
				listOfferingsMockFn: func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error) {
					return &servicemanager.ServiceOfferingResponseList{
						Items: []servicemanager.ServiceOfferingResponseObject{
							{
								Name: internal.Ptr("someOffering"),
								Id:   internal.Ptr("someID"),
							},
						},
					}, nil, nil
				},
				listPlansMockFn: func() (*servicemanager.ServicePlanResponseList, *http.Response, error) {
					return nil, nil, errors.New("plansApiError")
				},
			},
			wantErr: true,
		},
		{
			name: "empty response",
			args: args{
				listOfferingsMockFn: func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error) {
					return &servicemanager.ServiceOfferingResponseList{
						Items: []servicemanager.ServiceOfferingResponseObject{
							{
								Name: internal.Ptr("someOffering"),
								Id:   internal.Ptr("someID"),
							},
						},
					}, nil, nil
				},
				listPlansMockFn: func() (*servicemanager.ServicePlanResponseList, *http.Response, error) {
					return &servicemanager.ServicePlanResponseList{
						Items: []servicemanager.ServicePlanResponseObject{},
					}, nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				listOfferingsMockFn: func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error) {
					return &servicemanager.ServiceOfferingResponseList{
						Items: []servicemanager.ServiceOfferingResponseObject{
							{
								Name: internal.Ptr("someOffering"),
								Id:   internal.Ptr("someID"),
							},
						},
					}, nil, nil
				},
				listPlansMockFn: func() (*servicemanager.ServicePlanResponseList, *http.Response, error) {
					return &servicemanager.ServicePlanResponseList{
						Items: []servicemanager.ServicePlanResponseObject{
							{
								Name: internal.Ptr("somePlan"),
								Id:   internal.Ptr("somePlanID"),
							},
						},
					}, nil, nil
				},
			},
			wantErr: false,
			wantID:  "somePlanID",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smClient := &ServiceManagerClient{
				OfferingServiceFake{tc.args.listOfferingsMockFn},
				PlansServiceFake{listPlansMockFn: tc.args.listPlansMockFn},
			}
			planID, err := smClient.PlanIDByName(context.TODO(), "Not relevant, since mocked", "Not relevant, since mocked")

			if tc.wantErr != (err != nil) {
				t.Errorf("Unexpected error return; Expected error: %v, Returned: %v", tc.wantErr, err)
			}
			if tc.wantID != planID {
				t.Errorf("Unexpected returned PlanID; Expected: %s, Returned: %s", tc.wantID, planID)
			}

		})
	}
}

func TestNewCredsFromOperatorSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret map[string][]byte
		o      BindingCredentials
		err    error
	}{
		{
			name: "MissingAttributeError",
			secret: map[string][]byte{
				"clientid":     []byte("someClientId"),
				"clientsecret": []byte("someSecret"),
				"tokenurl":     []byte("https://valid.url"),
				"xsappname":    []byte("someXsAppName"),
			},
			err: errors.New(ErrInvalidSecretData),
		},
		{
			name: "SuccessfulMapping",
			secret: map[string][]byte{
				"clientid":     []byte("someClientId"),
				"clientsecret": []byte("someSecret"),
				"sm_url":       []byte("https://valid.url"),
				"tokenurl":     []byte("https://valid.url"),
				"xsappname":    []byte("someXsAppName"),
			},
			o: BindingCredentials{
				Clientid:     internal.Ptr("someClientId"),
				Clientsecret: internal.Ptr("someSecret"),
				SmUrl:        internal.Ptr("https://valid.url"),
				Url:          internal.Ptr("https://valid.url"),
				Xsappname:    internal.Ptr("someXsAppName"),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o, err := NewCredsFromOperatorSecret(tc.secret)
			if diff := cmp.Diff(tc.o, o); diff != "" {
				t.Errorf("\nNewBindingCredentialsFromSecretData(): -want, +got:\n%s\n", diff)
			}
			if diff := cmp.Diff(tc.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\nNewBindingCredentialsFromSecretData(): -want error, +got error:\n%s\n", diff)
			}

		})
	}
}

var _ servicemanager.ServiceOfferingsAPI = &OfferingServiceFake{}

type OfferingServiceFake struct {
	listOfferingsMockFn func() (*servicemanager.ServiceOfferingResponseList, *http.Response, error)
}

func (f OfferingServiceFake) GetServiceOfferingById(ctx context.Context, serviceOfferingID string) servicemanager.ApiGetServiceOfferingByIdRequest {
	panic("implement me")
}

func (f OfferingServiceFake) GetServiceOfferingByIdExecute(r servicemanager.ApiGetServiceOfferingByIdRequest) (*servicemanager.ServiceOfferingResponseObject, *http.Response, error) {
	panic("implement me")
}

func (f OfferingServiceFake) GetServiceOfferings(ctx context.Context) servicemanager.ApiGetServiceOfferingsRequest {
	return servicemanager.ApiGetServiceOfferingsRequest{ApiService: f}
}

func (f OfferingServiceFake) GetServiceOfferingsExecute(r servicemanager.ApiGetServiceOfferingsRequest) (*servicemanager.ServiceOfferingResponseList, *http.Response, error) {
	return f.listOfferingsMockFn()
}

var _ servicemanager.ServicePlansAPI = &PlansServiceFake{}

type PlansServiceFake struct {
	listPlansMockFn func() (*servicemanager.ServicePlanResponseList, *http.Response, error)
}

func (p PlansServiceFake) GetServicePlansByServiceId(ctx context.Context, servicePlanID string) servicemanager.ApiGetServicePlansByServiceIdRequest {
	panic("implement me")
}

func (p PlansServiceFake) GetServicePlansByServiceIdExecute(r servicemanager.ApiGetServicePlansByServiceIdRequest) (*servicemanager.ServicePlanResponseObject, *http.Response, error) {
	panic("implement me")
}

func (p PlansServiceFake) GetAllServicePlans(ctx context.Context) servicemanager.ApiGetAllServicePlansRequest {
	return servicemanager.ApiGetAllServicePlansRequest{ApiService: p}
}

func (p PlansServiceFake) GetAllServicePlansExecute(r servicemanager.ApiGetAllServicePlansRequest) (*servicemanager.ServicePlanResponseList, *http.Response, error) {
	return p.listPlansMockFn()
}
