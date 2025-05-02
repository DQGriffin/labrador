package types

// Assumed role types
type IamAssumedRolePolicy struct {
	Version   string
	Statement []IamAssumedRoleStatement
}

type IamAssumedRoleStatement struct {
	Effect    string
	Principal IamAssumedRolePrincipal
	Action    string
}

type IamAssumedRolePrincipal struct {
	Service []string
}

// Inline policy types
type IamInlinePolicy struct {
	Version   string
	Statement []IamInlinePolicyStatement
}

type IamInlinePolicyStatement struct {
	Effect   string
	Action   []string
	Resource []string
}
