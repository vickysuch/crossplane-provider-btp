package oidc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var token = "INLINE_TOKEN"

var kubeConfigInvalidTemplate = ">><>apiVersion: v1\nkind: Config\ncurrent-context: shoot--kyma-stage--c-51a159d\nclusters:\n- name: shoot--kyma-stage--c-51a159d\n  cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\ncontexts:\n- name: shoot--kyma-stage--c-51a159d\n  context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\nusers:\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - \"--oidc-issuer-url=xxx\"\n      - \"--oidc-client-id=xxx\"\n      - \"--oidc-extra-scope=email\"\n      - \"--oidc-extra-scope=openid\"\n      command: kubectl-oidc_login\n"
var kubeConfigNoUserTemplate = "apiVersion: v1\nkind: Config\ncurrent-context: shoot--kyma-stage--c-51a159d\nclusters:\n- name: shoot--kyma-stage--c-51a159d\n  cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\ncontexts:\n- name: shoot--kyma-stage--c-51a159d\n  context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\n"
var kubeConfigSingleUserTemplate = "apiVersion: v1\nkind: Config\ncurrent-context: shoot--kyma-stage--c-51a159d\nclusters:\n- name: shoot--kyma-stage--c-51a159d\n  cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\ncontexts:\n- name: shoot--kyma-stage--c-51a159d\n  context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\nusers:\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - \"--oidc-issuer-url=xxx\"\n      - \"--oidc-client-id=xxx\"\n      - \"--oidc-extra-scope=email\"\n      - \"--oidc-extra-scope=openid\"\n      command: kubectl-oidc_login\n      installHint: |\n        kubelogin plugin is required to proceed with authentication\n        # Homebrew (macOS and Linux)\n        brew install int128/kubelogin/kubelogin\n\n        # Krew (macOS, Linux, Windows and ARM)\n        kubectl krew install oidc-login\n\n        # Chocolatey (Windows)\n        choco install kubelogin"
var kubeConfigSingleUserTokenInline = fmt.Sprintf("apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\n  name: shoot--kyma-stage--c-51a159d\ncontexts:\n- context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\n  name: shoot--kyma-stage--c-51a159d\ncurrent-context: shoot--kyma-stage--c-51a159d\nkind: Config\nusers:\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - --oidc-issuer-url=xxx\n      - --oidc-client-id=xxx\n      - --oidc-extra-scope=email\n      - --oidc-extra-scope=openid\n      command: kubectl-oidc_login\n      installHint: |-\n        kubelogin plugin is required to proceed with authentication\n        # Homebrew (macOS and Linux)\n        brew install int128/kubelogin/kubelogin\n\n        # Krew (macOS, Linux, Windows and ARM)\n        kubectl krew install oidc-login\n\n        # Chocolatey (Windows)\n        choco install kubelogin\n    token: %s\n", token)
var kubeConfigSingleUserTokenReplace = fmt.Sprintf("apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\n  name: shoot--kyma-stage--c-51a159d\ncontexts:\n- context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\n  name: shoot--kyma-stage--c-51a159d\ncurrent-context: shoot--kyma-stage--c-51a159d\nkind: Config\nusers:\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    token: %s\n", token)

var kubeConfigMultiUserTemplate = "apiVersion: v1\nkind: Config\ncurrent-context: shoot--kyma-stage--c-51a159d\nclusters:\n  - name: shoot--kyma-stage--c-51a159d\n    cluster:\n      certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n      server: xxx\ncontexts:\n  - name: shoot--kyma-stage--c-51a159d\n    context:\n      cluster: shoot--kyma-stage--c-51a159d\n      user: shoot--kyma-stage--c-51a159d\nusers:\n  - name: shoot--kyma-stage--c-51a159d\n    user:\n      exec:\n        apiVersion: client.authentication.k8s.io/v1beta1\n        args:\n          - get-token\n          - \"--oidc-issuer-url=xxx\"\n          - \"--oidc-client-id=xxx\"\n          - \"--oidc-extra-scope=email\"\n          - \"--oidc-extra-scope=openid\"\n        command: kubectl-oidc_login\n  - name: shoot--kyma-stage--c-51a159d\n    user:\n      exec:\n        apiVersion: client.authentication.k8s.io/v1beta1\n        args:\n          - get-token\n          - \"--oidc-issuer-url=xxx\"\n          - \"--oidc-client-id=xxx\"\n          - \"--oidc-extra-scope=email\"\n          - \"--oidc-extra-scope=openid\"\n        command: kubectl-oidc_login\n  - name: shoot--kyma-stage--c-51a159d\n    user:\n      exec:\n        apiVersion: client.authentication.k8s.io/v1beta1\n        args:\n          - get-token\n          - \"--oidc-issuer-url=xxx\"\n          - \"--oidc-client-id=xxx\"\n          - \"--oidc-extra-scope=email\"\n          - \"--oidc-extra-scope=openid\"\n        command: kubectl-oidc_login"
var kubeConfigMultiUserTokenInline = fmt.Sprintf("apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\n  name: shoot--kyma-stage--c-51a159d\ncontexts:\n- context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\n  name: shoot--kyma-stage--c-51a159d\ncurrent-context: shoot--kyma-stage--c-51a159d\nkind: Config\nusers:\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - --oidc-issuer-url=xxx\n      - --oidc-client-id=xxx\n      - --oidc-extra-scope=email\n      - --oidc-extra-scope=openid\n      command: kubectl-oidc_login\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - --oidc-issuer-url=xxx\n      - --oidc-client-id=xxx\n      - --oidc-extra-scope=email\n      - --oidc-extra-scope=openid\n      command: kubectl-oidc_login\n    token: %s\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - --oidc-issuer-url=xxx\n      - --oidc-client-id=xxx\n      - --oidc-extra-scope=email\n      - --oidc-extra-scope=openid\n      command: kubectl-oidc_login\n", token)
var kubeConfigMultiUserTokenReplace = fmt.Sprintf("apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRUGdrSVZmUlRMa1dBeWhrVnRneU12REFOQmdrcWhraUc5dzBCQVFzRkFEQU4KTVFzd0NRWURWUVFERXdKallUQWVGdzB5TXpBeU1qSXdPVEk0TVRkYUZ3MHpNekF5TWpJd09USTRNVGRhTUEweApDekFKQmdOVkJBTVRBbU5oTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwQ2tpClNnS0xBWUQ5SklYZC9lMHk3NFF3aTFFSnYvQXpoRUNTVGlhb3J6R2RScy9UVDIrV0F3ME9YVitLRGxWd0hUMmwKNnFrQnAvMzQzcG9PSEVnQzJIRFFNdjR4MWdDTStzSWZCQjBHSEtLR1h0eW1tamcyc0M2dHg2eGFGYnB3cmR4ZgpURUk5T2hRWFVSaTZqcU40cjdVOWlRYWw0TlFQQkN4bXVwc0gvVmNueE1VOVNCWWNXc2ZUcCtiMGJJZ1hldUorCnZGdUJFOEovNytsZ2tqN2ZSWWlwcmlOQ2JDaFpvNEhqUmh2TzBwQTdvdk9JV3MrbUJlUEdVbHR5Nnp0NU9ETDgKNXZlVE5FbTlVOVBXYXl3Lzg2MTNndTBLTndaemM4c3RHV0V3Q0hGaWRqR1NJVGdXaXhBeE1rMExuemxFOURuSwpBK1ptR1RGekpXL3p6WmtKUFFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBYVl3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWdTd0NnUjFvSThySkNiQjNFcFBteXB3ZEhiTXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUQ1SURpelNSVG45NVVzTFY0T2hvT3QzZVFPNFVZaW4yNW03NGNneXpvWEV3K3hYVjFZRgpZNTR5LzY5SlNHTkJoK0dNMjBDQlQ4ZTVRSGFzZHAyTmhpTE1qR3VWSVNCcXJoQi9DZFB1OG5MZjczV21PbUFBClFyeXlsN2FYeXZyMGc3NVk0U1pwNkZrMzFrVDE0WFVqUFoyTDMxQTJyNzZWUmJKaXFNNG5CM1pYaHVBN1BsTC8KdFhWczdlY3dKdWRPYnNsMkVHQzZBRk9hbWxGTkx5Y3g5NENOYnJDTzIyZ1B4QlZXb1VHaWpQb0E2Z2F2WmV6TApoQU9ZSllFNDJTSm83a0FKb1VSU3VmU2dmZHFLY3lCejg5TmFJN1pGVlRzK1Jla0RxY0R3UXRybFJCNkFRMldrCm5wNWoyUlVxbWNOZTg5aHk1NFlRMm9vQ3JQV3VMeTliek9VPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\n    server: xxx\n  name: shoot--kyma-stage--c-51a159d\ncontexts:\n- context:\n    cluster: shoot--kyma-stage--c-51a159d\n    user: shoot--kyma-stage--c-51a159d\n  name: shoot--kyma-stage--c-51a159d\ncurrent-context: shoot--kyma-stage--c-51a159d\nkind: Config\nusers:\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - --oidc-issuer-url=xxx\n      - --oidc-client-id=xxx\n      - --oidc-extra-scope=email\n      - --oidc-extra-scope=openid\n      command: kubectl-oidc_login\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    token: %s\n- name: shoot--kyma-stage--c-51a159d\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - get-token\n      - --oidc-issuer-url=xxx\n      - --oidc-client-id=xxx\n      - --oidc-extra-scope=email\n      - --oidc-extra-scope=openid\n      command: kubectl-oidc_login\n", token)

func TestGenerate(t *testing.T) {

	tests := []struct {
		name string

		template   string
		token      string
		genOptions *GenerateConfig

		wantGenerated string
		wantErr       bool
	}{
		{
			name:          "Invalid Template, parsing error",
			template:      kubeConfigInvalidTemplate,
			token:         token,
			genOptions:    ConfigureGenerate().UserIndex(0).InjectInline(),
			wantGenerated: "",
			wantErr:       true,
		},
		{
			name:          "Missing User template",
			template:      kubeConfigNoUserTemplate,
			token:         token,
			genOptions:    ConfigureGenerate().UserIndex(0).InjectInline(),
			wantGenerated: "",
			wantErr:       true,
		},
		{
			name:          "SingleUserTemplate replace inline successful",
			template:      kubeConfigSingleUserTemplate,
			token:         token,
			genOptions:    ConfigureGenerate().UserIndex(0).InjectInline(),
			wantGenerated: kubeConfigSingleUserTokenInline,
			wantErr:       false,
		},
		{
			name:          "SingleUserTemplate replace whole user successful",
			template:      kubeConfigSingleUserTemplate,
			token:         token,
			genOptions:    ConfigureGenerate().UserIndex(0),
			wantGenerated: kubeConfigSingleUserTokenReplace,
			wantErr:       false,
		},
		{
			name:          "MultiUserTemplate replace inline successful",
			template:      kubeConfigMultiUserTemplate,
			token:         token,
			genOptions:    ConfigureGenerate().UserIndex(1).InjectInline(),
			wantGenerated: kubeConfigMultiUserTokenInline,
			wantErr:       false,
		},
		{
			name:          "MultiUserTemplate replace whole user successful",
			template:      kubeConfigMultiUserTemplate,
			token:         token,
			genOptions:    ConfigureGenerate().UserIndex(1),
			wantGenerated: kubeConfigMultiUserTokenReplace,
			wantErr:       false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			creator := testKubeConfigCreator()
			generate, err := creator.Generate([]byte(tc.template), []byte(tc.token), tc.genOptions)
			if (tc.wantErr && err == nil) || (!tc.wantErr && err != nil) {
				t.Errorf("\n%s\ne.Generate(...): unexpected error behaviour, got\n%s", tc.name, err)
			}
			if tc.wantGenerated != "" {
				if diff := cmp.Diff([]byte(tc.wantGenerated), generate.GeneratedKubeConfig); diff != "" {
					t.Errorf("\n%s\ne.Generate(...): -want, +got:\n--------\n%s\n++++++++\n%s\n", tc.name, tc.wantGenerated, generate.GeneratedKubeConfig)
				}
			}
		})
	}
}
func TestIsUpToDate(t *testing.T) {
	configBytes := []byte(kubeConfigSingleUserTemplate)
	tokenBytes := []byte(token)

	configHash, tokenHash := generateHashes(t, configBytes, tokenBytes)
	// no configured previous hash
	boolTestCase(t,
		false,
		testKubeConfigCreator().WithHashes(nil, nil),
		configBytes, tokenBytes,
	)
	// token hash different -> not matching
	boolTestCase(t,
		false,
		testKubeConfigCreator().WithHashes(configHash, configHash),
		configBytes, tokenBytes,
	)
	// kubeconfig hash different -> not matching
	boolTestCase(t,
		false,
		testKubeConfigCreator().WithHashes(tokenHash, tokenHash),
		configBytes, tokenBytes,
	)
	// same hashes -> matching
	boolTestCase(t,
		true,
		testKubeConfigCreator().WithHashes(configHash, tokenHash),
		configBytes, tokenBytes,
	)
}

func testKubeConfigCreator() KubeConfigClient {
	return NewKubeConfigCreator([]byte{}, []byte{})
}

func boolTestCase(t *testing.T, expected bool, creator KubeConfigClient, bytes []byte, tokenBytes []byte) {
	assert.Equal(t, expected, creator.IsUpToDate(bytes, tokenBytes))
}
func generateHashes(t *testing.T, configBytes []byte, tokenBytes []byte) ([]byte, []byte) {
	generate, err := testKubeConfigCreator().Generate(configBytes, tokenBytes, ConfigureGenerate())
	if err != nil {
		assert.Nil(t, err)
	}
	return generate.SourceKubeConfigHash, generate.SourceTokenHash
}
