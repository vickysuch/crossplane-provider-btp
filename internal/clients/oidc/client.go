package oidc

import (
	"context"

	"github.com/int128/kubelogin/pkg/oidc"
)

type LoginPerformer interface {
	DoLogin(ctx context.Context) (*oidc.TokenSet, error)
	IsExpired(idToken string) bool
	Refresh(ctx context.Context, refreshToken string) (*oidc.TokenSet, error)
}

type KubeConfigClient interface {
	WithHashes(kubeConfigHash []byte, tokenHash []byte) KubeConfigClient
	IsUpToDate(kubeConfig []byte, token []byte) bool
	Generate(kubeConfig []byte, token []byte, config *GenerateConfig) (GenerateResult, error)
}

type GenerateResult struct {
	SourceKubeConfigHash []byte
	SourceTokenHash      []byte
	GeneratedKubeConfig  []byte
	ServerUrl            string
}

// Configuration to control process of generation kubeconfig from template
type GenerateConfig struct {
	userIndex int
	setInline bool
}

// Entry point for building configuration for generate process
func ConfigureGenerate() *GenerateConfig {
	return &GenerateConfig{
		userIndex: 0,
		setInline: false,
	}
}

// User to injec the user to, starting with index 0
func (gConfig *GenerateConfig) UserIndex(i int) *GenerateConfig {
	gConfig.userIndex = i
	return gConfig
}

// If true replaces only token, but leaves rest of user untouched
func (gConfig *GenerateConfig) InjectInline() *GenerateConfig {
	gConfig.setInline = true
	return gConfig
}
