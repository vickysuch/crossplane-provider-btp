package v1alpha1

const (
	// SubaccountOperatorLabel Unique identifier of this operator,
	//workaround as the label query for cis local does not work with dots (.)
	SubaccountOperatorLabel = "orchestrate.cloud.sap/subaccount-operator"
	SMLabel                 = "orchestrate_cloud_sap_subaccount_operator"
)
