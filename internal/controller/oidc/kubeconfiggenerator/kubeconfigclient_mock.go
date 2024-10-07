package kubeconfiggenerator

import (
	"github.com/pkg/errors"
	"github.com/sap/crossplane-provider-btp/internal/clients/oidc"
)

var _ oidc.KubeConfigClient = &KubeConfigClientMock{}

var errGenerate = errors.New("Generate error")

var (
	MockedKubeConfigHash = []byte("KUBECONFIGHASH")
	MockedTokenHash      = []byte("TOKENHASH")
)

type KubeConfigClientMock struct {
	UpToDate         bool
	GeneratedContent string
	ServerUrl        string
}

func (k *KubeConfigClientMock) WithHashes(kubeConfigHash []byte, tokenHash []byte) oidc.KubeConfigClient {
	return k
}

func (k *KubeConfigClientMock) IsUpToDate(kubeConfig []byte, token []byte) bool {
	return k.UpToDate
}

func (k *KubeConfigClientMock) Generate(kubeConfig []byte, token []byte, config *oidc.GenerateConfig) (oidc.GenerateResult, error) {
	if k.GeneratedContent == "" {
		return oidc.GenerateResult{}, errGenerate
	}
	return oidc.GenerateResult{
		SourceKubeConfigHash: MockedKubeConfigHash,
		SourceTokenHash:      MockedTokenHash,
		GeneratedKubeConfig:  []byte(k.GeneratedContent),
		ServerUrl:            k.ServerUrl,
	}, nil
}
