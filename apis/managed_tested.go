package apis

import "github.com/crossplane/crossplane-runtime/pkg/resource"

type ManagedTested interface {
	resource.Managed
	SetExternalID(newID string)
	GetExternalID() string
}
