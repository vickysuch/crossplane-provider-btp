package rolecollectionassignment

import (
	"context"
)

type RoleAssignerMock struct {
	hasRole bool
	err     error

	CalledIdentifier *string
}

func (u *RoleAssignerMock) RevokeRole(ctx context.Context, origin, identifier, rolecollection string) error {
	u.CalledIdentifier = &identifier
	return u.err
}

func (u *RoleAssignerMock) AssignRole(ctx context.Context, origin, identifier, rolecollection string) error {
	u.CalledIdentifier = &identifier
	return u.err
}

func (u *RoleAssignerMock) HasRole(ctx context.Context, origin, identifier, roleCollection string) (bool, error) {
	u.CalledIdentifier = &identifier
	return u.hasRole, u.err
}

var _ RoleAssigner = &RoleAssignerMock{}
