package btp

import "github.com/pkg/errors"

func NewBTPClient(cisSecretData []byte, serviceAccountSecretData []byte) (*Client, error) {

	accountsServiceClient, err := ServiceClientFromSecret(cisSecretData, serviceAccountSecretData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get BTP accounts service client.")
	}
	return &accountsServiceClient, nil
}
