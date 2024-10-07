package oidc

import (
	"bytes"
	"crypto/sha256"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	"github.com/sap/crossplane-provider-btp/internal"
)

var (
	errMalformedKubeConfig = "kubeconfig does not match expected format"
	errMissConfiguration   = "User Index does not exist in provided kube config template"
	errGenerateKubeConfig  = "error while generating kubeconfig from template"
)

var _ KubeConfigClient = &KubeConfigCreator{}

func NewKubeConfigCreator(tokenHash []byte, kubeConfigHash []byte) KubeConfigClient {
	k := &KubeConfigCreator{}
	return k.WithHashes(kubeConfigHash, tokenHash)
}

type KubeConfigCreator struct {
	tokenHash      []byte
	kubeConfigHash []byte
}

func (k *KubeConfigCreator) WithHashes(kubeConfigHash []byte, tokenHash []byte) KubeConfigClient {
	k.kubeConfigHash = kubeConfigHash
	k.tokenHash = tokenHash
	return k
}

func (k *KubeConfigCreator) IsUpToDate(kubeConfig []byte, token []byte) bool {
	if k.tokenHash == nil || k.kubeConfigHash == nil {
		return false
	}
	return bytes.Equal(hash(kubeConfig), k.kubeConfigHash) && bytes.Equal(hash(token), k.tokenHash)
}

func (k *KubeConfigCreator) Generate(kubeConfig []byte, token []byte, config *GenerateConfig) (GenerateResult, error) {
	if config == nil {
		return GenerateResult{}, errors.New("Generate process can't be run without a configuration")
	}
	generate, err := injectToken(kubeConfig, token, config)
	if err != nil {
		return GenerateResult{}, errors.Wrap(err, errGenerateKubeConfig)
	}

	serverUrl, _ := internal.ParseConnectionDetailsFromKubeYaml(generate)

	result := GenerateResult{
		SourceKubeConfigHash: hash(kubeConfig),
		SourceTokenHash:      hash(token),
		GeneratedKubeConfig:  generate,
		ServerUrl:            serverUrl,
	}
	return result, nil
}

func injectToken(kubeConfig []byte, token []byte, config *GenerateConfig) ([]byte, error) {
	var kConfig map[string]interface{}
	marshalError := yaml.Unmarshal(kubeConfig, &kConfig)
	if marshalError != nil {
		return nil, marshalError
	}
	user, err := parseUserUnstructured(kConfig, config.userIndex)
	if err != nil {
		return nil, err
	}
	user["token"] = string(token)
	if !config.setInline {
		unsetOtherKeys(user, "token")
	}
	return yaml.Marshal(kConfig)
}

func parseUserUnstructured(config map[string]interface{}, index int) (map[string]interface{}, error) {
	uList := config["users"]
	uListTyped, ok := uList.([]interface{})
	if !ok {
		return nil, errors.New(errMalformedKubeConfig)
	}
	//TODO: can use another test case
	if len(uListTyped) <= index {
		return nil, errors.New(errMissConfiguration)
	}
	userWrap := uListTyped[index]
	userWrapTyped, ok := userWrap.(map[string]interface{})
	if !ok {
		return nil, errors.New(errMalformedKubeConfig)
	}
	user := userWrapTyped["user"]
	userTyped, ok := user.(map[string]interface{})
	if !ok {
		return nil, errors.New(errMalformedKubeConfig)
	}
	return userTyped, nil
}

func hash(in []byte) []byte {
	sum256 := sha256.Sum256(in)
	return sum256[:]
}

func unsetOtherKeys(m map[string]interface{}, keepKey string) {
	for k := range m {
		if k != keepKey {
			delete(m, k)
		}
	}
}
