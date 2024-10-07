package environments

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sap/crossplane-provider-btp/apis/environment/v1alpha1"
	"github.com/sap/crossplane-provider-btp/btp"
)

func TestNeedsUpdate(t *testing.T) {
	techUser := "techUser@sap.com"

	tests := []struct {
		name      string
		cr        v1alpha1.CloudFoundryEnvironment
		credsUser string
		want      bool
	}{
		{
			name:      "NoUpdateIgnoreTechUser",
			cr:        crWithManagers([]string{"1@sap.com", "2@sap.com", techUser}, []string{"1@sap.com", "2@sap.com", techUser}),
			credsUser: techUser,
			want:      false,
		},
		{
			name:      "NoUpdateIgnoreAddingTechUser",
			cr:        crWithManagers([]string{"1@sap.com", "2@sap.com"}, []string{"1@sap.com", "2@sap.com", techUser}),
			credsUser: techUser,
			want:      false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uua := CloudFoundryOrganization{
				btp: btp.Client{Credential: &btp.Credentials{UserCredential: &btp.UserCredential{Email: tc.credsUser}}},
			}
			needsUpdate := uua.NeedsUpdate(tc.cr)

			if diff := cmp.Diff(tc.want, needsUpdate); diff != "" {
				t.Errorf("\ne.NeedsUpdate(...): -want, +got:\n%s\n", diff)
			}

		})
	}
}

func TestManagerDiff(t *testing.T) {
	techUser := "techUser@sap.com"
	techUserUpper := strings.ToUpper(techUser)

	type want struct {
		add    []v1alpha1.User
		remove []v1alpha1.User
	}

	tests := []struct {
		name      string
		cr        v1alpha1.CloudFoundryEnvironment
		credsUser string
		want      want
	}{
		{
			name:      "Equal",
			cr:        crWithManagers([]string{"1@sap.com", "2@sap.com"}, []string{"1@sap.com", "2@sap.com"}),
			credsUser: techUserUpper,
			want:      want{add: []v1alpha1.User{}, remove: []v1alpha1.User{}},
		},
		{
			name:      "EqualWithTechUser",
			cr:        crWithManagers([]string{"1@sap.com", "2@sap.com", techUser}, []string{"1@sap.com", "2@sap.com", techUser}),
			credsUser: techUserUpper,
			want:      want{add: []v1alpha1.User{}, remove: []v1alpha1.User{}},
		},
		{
			name:      "NotAddTechUser",
			cr:        crWithManagers([]string{"1@sap.com", "2@sap.com", techUser}, []string{"1@sap.com", "2@sap.com"}),
			credsUser: techUserUpper,
			want:      want{add: []v1alpha1.User{}, remove: []v1alpha1.User{}},
		},
		{
			name:      "NotRemoveTechUser",
			cr:        crWithManagers([]string{"1@sap.com", "2@sap.com"}, []string{"1@sap.com", "2@sap.com", techUser}),
			credsUser: techUserUpper,
			want:      want{add: []v1alpha1.User{}, remove: []v1alpha1.User{}},
		},
		{
			name:      "AddAndRemove",
			cr:        crWithManagers([]string{"2@sap.com", "3@sap.com"}, []string{"1@sap.com", "2@sap.com"}),
			credsUser: techUserUpper,
			want:      want{add: toUserSlice([]string{"3@sap.com"}), remove: toUserSlice([]string{"1@sap.com"})},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uua := CloudFoundryOrganization{
				btp: btp.Client{Credential: &btp.Credentials{UserCredential: &btp.UserCredential{Email: tc.credsUser}}},
			}
			add, remove := uua.managerDiff(tc.cr)

			if diff := cmp.Diff(tc.want.add, add); diff != "" {
				t.Errorf("\ne.managerDiff(...): -want add, +got add:\n%s\n", diff)
			}
			if diff := cmp.Diff(tc.want.remove, remove); diff != "" {
				t.Errorf("\ne.managerDiff(...): -want remove, +got remove:\n%s\n", diff)
			}
		})
	}
}

func crWithManagers(specManagers []string, statusManagers []string) v1alpha1.CloudFoundryEnvironment {

	return v1alpha1.CloudFoundryEnvironment{
		Spec: v1alpha1.CfEnvironmentSpec{
			ForProvider: v1alpha1.CfEnvironmentParameters{
				Managers: specManagers,
			},
		},
		Status: v1alpha1.EnvironmentStatus{
			AtProvider: v1alpha1.CfEnvironmentObservation{
				Managers: toUserSlice(statusManagers),
			},
		},
	}
}

func toUserSlice(ss []string) []v1alpha1.User {
	us := make([]v1alpha1.User, 0)
	for _, s := range ss {
		us = append(us, v1alpha1.User{Username: s, Origin: "sap.ids"})
	}

	return us
}
