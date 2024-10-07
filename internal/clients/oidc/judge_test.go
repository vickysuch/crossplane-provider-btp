package oidc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sap/crossplane-provider-btp/apis/oidc/v1alpha1"
	"github.com/sap/crossplane-provider-btp/internal"
	"github.com/sap/crossplane-provider-btp/internal/testutils"
)

const anHourInPast = time.Hour * -1
const anHourInFuture = time.Hour * 1
const halfHour = time.Minute * 30
const quarterHour = time.Minute * 15

func TestJwtJudge_EstimateRotationStart(t *testing.T) {
	type args struct {
		idToken string
	}
	const halfHour = time.Minute * 30
	noError := func(t assert.TestingT, err error, i ...interface{}) bool {
		return err != nil
	}
	noIATError :=
		func(t assert.TestingT, err error, i ...interface{}) bool {
			return err.Error() == errCouldNotExtractIssuedAt
		}
	noExpError :=
		func(t assert.TestingT, err error, i ...interface{}) bool {
			return err.Error() == errCouldNotExtractExpiry
		}
	negativeDurationError :=
		func(t assert.TestingT, err error, i ...interface{}) bool {
			return err.Error() == errDurationNegative
		}
	tests := []struct {
		name    string
		args    args
		want    *time.Duration
		wantErr assert.ErrorAssertionFunc
	}{
		{name: "30mins", args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.IssuedAt(0), testutils.ExpiresAt(halfHour*3))}, want: toDurationPtr(halfHour), wantErr: noError},
		{name: "60mins", args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.IssuedAt(time.Hour*2), testutils.ExpiresAt(time.Hour*3))}, want: toDurationPtr(time.Hour * 1), wantErr: noError},
		{name: "No issued at, returns err", args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.ExpiresAt(time.Hour*3))}, want: nil, wantErr: noIATError},
		{name: "No exp at at, returns err", args: args{idToken: testutils.JwtToken(testutils.Epoch)}, want: nil, wantErr: noExpError},
		{name: "issued at before expires, returns err", args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.IssuedAt(halfHour), testutils.ExpiresAt(0))}, want: nil, wantErr: negativeDurationError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			judge, _ := NewJwtJudge(tt.args.idToken)
			got, err := judge.EstimateRotationDuration()
			if !tt.wantErr(t, err, fmt.Sprintf("EstimateRotationDuration(%v)", tt.args.idToken)) {
				return
			}
			assert.Equalf(t, tt.want, got, "EstimateRotationDuration(%v)", tt.args.idToken)
		})
	}
}

func toDurationPtr(num time.Duration) *time.Duration {
	dur := num
	return &dur
}
func TestJwtJudge_IsInRenewPeriod(t *testing.T) {

	type args struct {
		idToken          string
		renewPeriodStart time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{name: "Valid token - no time based claims, renew in future - no renew needed", args: args{idToken: testutils.JwtToken(testutils.Now), renewPeriodStart: anHourInFuture}, want: false},
		{name: "Valid token - issued in the past, renew in future- returns false", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.IssuedAt(anHourInPast)), renewPeriodStart: halfHour}, want: false},
		{name: "Valid token - nbf in the past, renew in future- returns false", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.NotBefore(anHourInPast)), renewPeriodStart: halfHour}, want: false},
		{name: "Valid token - expiry in future, renew in future - no renew needed", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.ExpiresAt(anHourInFuture)), renewPeriodStart: halfHour}, want: false},

		{name: "Invalid token - valid in future - needs renew", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.NotBefore(anHourInFuture))}, want: true},
		{name: "Invalid token - issued in future - needs renew", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.IssuedAt(anHourInFuture))}, want: true},
		{name: "Invalid token - expired an hour ago - needs renew", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.ExpiresAt(anHourInPast))}, want: true},

		{name: "Will expire in 15mins, renew started half an hour ago - needs renew", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.ExpiresAt(quarterHour)), renewPeriodStart: halfHour}, want: true},
		{name: "Will expire in 15mins, iat in the past, renew started half an hour ago - needs renew", args: args{idToken: testutils.JwtToken(testutils.Now, testutils.IssuedAt(anHourInPast), testutils.ExpiresAt(quarterHour)), renewPeriodStart: halfHour}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			judge, _ := NewJwtJudge(tt.args.idToken)
			isInRenewPeriod := judge.IsInRenewPeriod(tt.args.renewPeriodStart)
			assert.Equalf(t, tt.want, isInRenewPeriod, "HasReachedRenewPeriod(%v)", tt.args.idToken)
		})
	}
}

func TestJwtJudge_Status(t *testing.T) {
	type args struct {
		idToken string
	}
	tests := []struct {
		name string
		args args
		want v1alpha1.JwtStatus
	}{
		{
			name: "No claims in jwt",
			args: args{idToken: testutils.JwtToken(testutils.Epoch)},
			want: v1alpha1.JwtStatus{
				RotationStrategy: internal.Ptr(v1alpha1.RotationStrategyDynamic),
			},
		},
		{
			name: "Only iat in jwt",
			args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.IssuedAt(0))},
			want: v1alpha1.JwtStatus{
				RotationStrategy: internal.Ptr(v1alpha1.RotationStrategyDynamic),
				IssuedAt:         internal.Ptr(metav1.Time{Time: testutils.Epoch()}),
			},
		},
		{
			name: "Only iss in jwt",
			args: args{idToken: testutils.JwtTokenWithIssuer("https://intern.test", testutils.Epoch)},
			want: v1alpha1.JwtStatus{
				Issuer:           internal.Ptr("https://intern.test"),
				RotationStrategy: internal.Ptr(v1alpha1.RotationStrategyDynamic),
			},
		},
		{
			name: "Only exp in jwt",
			args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.ExpiresAt(0))},
			want: v1alpha1.JwtStatus{
				RotationStrategy: internal.Ptr(v1alpha1.RotationStrategyDynamic),
				ExpiresAt:        internal.Ptr(metav1.Time{Time: testutils.Epoch()}),
			},
		},
		{
			name: "iat and exp in jwt, expect rotation properties set",
			args: args{idToken: testutils.JwtToken(testutils.Epoch, testutils.IssuedAt(0), testutils.ExpiresAt(time.Hour*3))},
			want: v1alpha1.JwtStatus{
				RotationStrategy:  internal.Ptr(v1alpha1.RotationStrategyDynamic),
				IssuedAt:          internal.Ptr(metav1.Time{Time: testutils.Epoch()}),
				ExpiresAt:         internal.Ptr(metav1.Time{Time: testutils.Epoch().Add(time.Hour * 3)}),
				RotationDuration:  internal.Ptr(metav1.Duration{Duration: time.Hour * 1}),
				RotationNotBefore: internal.Ptr(metav1.Time{Time: testutils.Epoch().Add(time.Hour * 2)}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			judge, _ := NewJwtJudge(tt.args.idToken)
			assert.Equalf(t, tt.want, judge.Status(), "Status()")
		})
	}
}
