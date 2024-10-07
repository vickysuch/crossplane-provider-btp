package rolecollection

import (
	"context"

	"github.com/sap/crossplane-provider-btp/apis/security/v1alpha1"
)

// RoleMaintainerMock is a mock implementation of RoleCollectionMaintainer interface
// returns stubed values and records called identifier to most methods
type RoleMaintainerMock struct {
	generateObservation v1alpha1.RoleCollectionObservation
	needsCreation       bool
	needsUpdate         bool
	err                 error
	// for verification
	CalledIdentifier string
}

var _ RoleCollectionMaintainer = &RoleMaintainerMock{}

func (r *RoleMaintainerMock) GenerateObservation(ctx context.Context, roleCollectionName string) (v1alpha1.RoleCollectionObservation, error) {
	r.CalledIdentifier = roleCollectionName
	return r.generateObservation, r.err
}

func (r *RoleMaintainerMock) NeedsCreation(collection v1alpha1.RoleCollectionObservation) bool {
	return r.needsCreation
}

func (r *RoleMaintainerMock) NeedsUpdate(params v1alpha1.RoleCollectionParameters, observation v1alpha1.RoleCollectionObservation) bool {
	return r.needsUpdate
}

func (r *RoleMaintainerMock) Create(ctx context.Context, params v1alpha1.RoleCollectionParameters) (string, error) {
	return r.CalledIdentifier, r.err
}

func (r *RoleMaintainerMock) Update(ctx context.Context, roleCollectionName string, params v1alpha1.RoleCollectionParameters, obs v1alpha1.RoleCollectionObservation) error {
	r.CalledIdentifier = roleCollectionName
	return r.err
}

func (r *RoleMaintainerMock) Delete(ctx context.Context, roleCollectionName string) error {
	r.CalledIdentifier = roleCollectionName
	return r.err
}
