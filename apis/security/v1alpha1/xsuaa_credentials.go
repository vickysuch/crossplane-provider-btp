package v1alpha1

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var ErrInvalidXsuaaCredentials = errors.New("invalid xsuaa api credentials")

// XsuaaBinding defines the json structure stored in secret to configure xsuaa api client
type XsuaaBinding struct {
	ClientId     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	TokenURL     string `json:"tokenurl"`
	ApiUrl       string `json:"apiurl"`
}

func ReadXsuaaCredentials(creds []byte) (XsuaaBinding, error) {
	var binding = XsuaaBinding{}
	if err := json.Unmarshal(creds, &binding); err != nil {
		return binding, ErrInvalidXsuaaCredentials
	}
	if binding.ClientId == "" || binding.ClientSecret == "" || binding.TokenURL == "" || binding.ApiUrl == "" {
		return binding, ErrInvalidXsuaaCredentials
	}
	return binding, nil
}
