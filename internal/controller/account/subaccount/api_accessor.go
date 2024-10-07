package subaccount

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sap/crossplane-provider-btp/btp"
	accountclient "github.com/sap/crossplane-provider-btp/internal/openapi_clients/btp-accounts-service-api-go/pkg"
)

// AccountsApiAccessor abstraction to handle API operations by coordinating to generated api client
type AccountsApiAccessor interface {
	MoveSubaccount(ctx context.Context, subaccountGuid string, targetId string) error
	UpdateSubaccount(ctx context.Context, subaccountGuid string, payload accountclient.UpdateSubaccountRequestPayload) error
}

type AccountsClient struct {
	btp btp.Client
}

func (a *AccountsClient) UpdateSubaccount(ctx context.Context, subaccountGuid string, payload accountclient.UpdateSubaccountRequestPayload) error {
	_, _, err := a.btp.AccountsServiceClient.SubaccountOperationsAPI.
		UpdateSubaccount(ctx, subaccountGuid).
		UpdateSubaccountRequestPayload(payload).
		Execute()
	return err
}

func (a *AccountsClient) MoveSubaccount(ctx context.Context, subaccountGuid string, targetId string) error {
	if targetId == "" {
		return errors.New("targetId must be set for move subaccount api call")
	}
	_, _, err := a.btp.AccountsServiceClient.SubaccountOperationsAPI.
		MoveSubaccount(ctx, subaccountGuid).
		MoveSubaccountRequestPayload(
			accountclient.MoveSubaccountRequestPayload{TargetAccountGUID: targetId}).
		Execute()
	return err
}

var _ AccountsApiAccessor = &AccountsClient{}
