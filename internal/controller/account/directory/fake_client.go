package directory

import (
	"context"

	"github.com/sap/crossplane-provider-btp/apis/account/v1alpha1"
	"github.com/sap/crossplane-provider-btp/internal/clients/directory"
)

type MockClient struct {
	needsCreation    bool
	needsCreationErr error

	needsUpdate bool

	createErr    error
	createResult v1alpha1.Directory

	updateErr error

	deleteErr error

	syncErr error

	available bool
}

func (d MockClient) IsAvailable() bool {
	return d.available
}

func (d MockClient) SyncStatus(ctx context.Context) error {
	return d.syncErr
}

func (d MockClient) UpdateDirectory(ctx context.Context) (*v1alpha1.Directory, error) {
	return nil, d.updateErr
}

func (d MockClient) CreateDirectory(ctx context.Context) (*v1alpha1.Directory, error) {
	return &d.createResult, d.createErr
}

func (d MockClient) NeedsUpdate(ctx context.Context) (bool, error) {
	return d.needsUpdate, nil
}

func (d MockClient) NeedsCreation(ctx context.Context) (bool, error) {
	return d.needsCreation, d.needsCreationErr
}

func (d MockClient) DeleteDirectory(ctx context.Context) error {
	return d.deleteErr
}

var _ directory.DirectoryClientI = &MockClient{}
