package certbasedoidclogin

import (
	"context"
	oidc2 "github.com/int128/kubelogin/pkg/oidc"
	"github.com/pkg/errors"
	"github.com/sap/crossplane-provider-btp/internal/clients/oidc"
)

var (
	errMockedLogin   = errors.New("mocked login error")
	errMockedRefresh = errors.New("refresh error")
)

type CertLoginMock struct {
	TokenSet *oidc2.TokenSet
	Expired  bool
}

var _ oidc.LoginPerformer = &CertLoginMock{}

func (cLMock *CertLoginMock) DoLogin(ctx context.Context) (*oidc2.TokenSet, error) {
	if cLMock.TokenSet == nil {
		return nil, errMockedLogin
	}
	return cLMock.TokenSet, nil
}

func (cLMock *CertLoginMock) IsExpired(idToken string) bool {
	return cLMock.Expired
}

func (cLMock *CertLoginMock) Refresh(ctx context.Context, refreshToken string) (*oidc2.TokenSet, error) {
	if cLMock.TokenSet == nil {
		return nil, errMockedRefresh
	}
	return cLMock.TokenSet, nil
}
