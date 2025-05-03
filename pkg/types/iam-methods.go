package types

import (
	"strings"

	"github.com/DQGriffin/labrador/internal/cli/styles"
	"github.com/DQGriffin/labrador/pkg/helpers"
	"github.com/charmbracelet/lipgloss/tree"
)

func (config IamRoleConfig) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	for _, role := range config.Roles {
		childNodes := role.ToTreeNodes(verbose)
		nodes = append(nodes, childNodes...)
	}

	return nodes
}

func (role IamRoleSettings) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(*role.Name))

	if verbose {
		node.Child(styles.Primary.Render("Ref:          ") + styles.Secondary.Render(helpers.PtrOrDefault(role.Ref, "[not set]")))
		node.Child(styles.Primary.Render("Description:  ") + styles.Secondary.Render(helpers.PtrOrDefault(role.Description, "[not set]")))
		arnNodes := generateIamRolePolicyArnNodes(&role.PolicyArns)
		for i := range arnNodes {
			node.Child(arnNodes[i])
		}

		inlinePolicyNodes := generateIamRoleInlinePolicyNodes(&role.InlinePolicies)
		for i := range inlinePolicyNodes {
			node.Child(inlinePolicyNodes[i])
		}

		trustPolicyNodes := generateIamRoleTrustPolicyNodes(*role.TrustPolicy)
		for i := range trustPolicyNodes {
			node.Child(trustPolicyNodes[i])
		}
	}

	nodes = append(nodes, node)
	return nodes
}

func generateIamRolePolicyArnNodes(policyArns *[]string) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render("Policy ARNs"))
	for _, arn := range *policyArns {
		node.Child(styles.Primary.Render("ARN:  ") + styles.Secondary.Render(arn))
	}

	nodes = append(nodes, node)
	return nodes
}

func generateIamRoleInlinePolicyNodes(policies *[]IamInlinePolicy) []*tree.Tree {
	var nodes []*tree.Tree
	root := tree.New().Root(styles.Primary.Render("Inline Policies"))
	for _, policy := range *policies {
		node := tree.New().Root(styles.Primary.Render(policy.Name))
		node.Child(styles.Primary.Render("Effect:     ") + styles.Secondary.Render(*policy.Effect))
		node.Child(styles.Primary.Render("Actions:    ") + styles.Secondary.Render(strings.Join(policy.Actions, ", ")))
		node.Child(styles.Primary.Render("Resources:  ") + styles.Secondary.Render(strings.Join(policy.Resources, ", ")))
		root.Child(node)
	}

	nodes = append(nodes, root)
	return nodes
}

func generateIamRoleTrustPolicyNodes(trustPolicy IamTrustPolicy) []*tree.Tree {
	var nodes []*tree.Tree
	root := tree.New().Root(styles.Primary.Render("Trust Policy"))

	if trustPolicy.FilePath != nil {
		root.Child(styles.Primary.Render("File:      ") + styles.Secondary.Render(*trustPolicy.FilePath))
	} else {
		root.Child(styles.Primary.Render("Services:      ") + styles.Secondary.Render(strings.Join(trustPolicy.Principals.Services, ", ")))
		root.Child(styles.Primary.Render("AWS Accounts:  ") + styles.Secondary.Render(strings.Join(trustPolicy.Principals.AwsAccounts, ", ")))
	}

	nodes = append(nodes, root)
	return nodes
}
