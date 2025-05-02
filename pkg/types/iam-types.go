package types

type IamRoleConfig struct {
	Defaults *IamRoleSettings  `json:"defaults,omitempty"`
	Roles    []IamRoleSettings `json:"roles,omitempty"`
}

type IamRoleSettings struct {
	Name           *string           `json:"name,omitempty"`
	Ref            *string           `json:"ref,omitempty"`
	Description    *string           `json:"description,omitempty"`
	TrustPolicy    *IamTrustPolicy   `json:"trustPolicy,omitempty"`
	PolicyArns     []string          `json:"policyArns,omitempty"`
	InlinePolicies []IamInlinePolicy `json:"inlinePolicies,omitempty"`
}

type IamTrustPolicy struct {
	Principals *IamPrincipals `json:"principals,omitempty"`
	FilePath   *string        `json:"file,omitempty"`
}

type IamPrincipals struct {
	Services    []string `json:"services,omitempty"`
	AwsAccounts []string `json:"aws,omitempty"`
}

type IamInlinePolicy struct {
	Name      string   `json:"name,omitempty"`
	Actions   []string `json:"actions,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Effect    *string  `json:"effect,omitempty"`
	FilePath  *string  `json:"file,omitempty"`
}

type IamPermission struct {
	Actions   []string `json:"actions,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Effect    string   `json:"effect,omitempty"`
	FilePath  *string  `json:"file,omitempty"`
}
