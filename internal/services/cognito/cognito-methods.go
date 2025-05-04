package cognito

import (
	"fmt"
	"strings"

	"github.com/DQGriffin/labrador/internal/cli/styles"
	"github.com/charmbracelet/lipgloss/tree"
)

func (config CognitoConfig) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	for _, pool := range config.Pools {
		childNodes := pool.ToTreeNodes(verbose)
		nodes = append(nodes, childNodes...)
	}

	return nodes
}

func (cognito CognitoSettings) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(*cognito.ApplicationName))

	if verbose {
		node.Child(styles.Primary.Render("Sign-up Atrributes:   ") + styles.Secondary.Render(cognito.JoinedSignUpAttributes()))
		node.Child(styles.Primary.Render("Sign-in Identifiers:  ") + styles.Secondary.Render(cognito.JoinedSignInIdentifiers()))

		passwordPolicyNodes := generatePasswordPolicytTreeNodes(cognito.PasswordRequirements)
		for i := range passwordPolicyNodes {
			node.Child(passwordPolicyNodes[i])
		}

		appClientNodes := generateAppClientTreeNodes(cognito.AppClients)
		for i := range appClientNodes {
			node.Child(appClientNodes[i])
		}
	}

	nodes = append(nodes, node)
	return nodes
}

func generateAppClientTreeNodes(clients *[]CognitoAppClient) []*tree.Tree {
	var nodes []*tree.Tree
	root := tree.New().Root(styles.Primary.Render("App Clients"))

	for _, client := range *clients {
		node := tree.New().Root(styles.Primary.Render(client.Name))
		node.Child(styles.Primary.Render("Client Type:  ") + styles.Secondary.Render(client.ClientType))
		node.Child(styles.Primary.Render("Return URLs:  ") + styles.Secondary.Render(client.JoinedReturnUrls()))

		root.Child(node)
	}

	nodes = append(nodes, root)
	return nodes
}

func generatePasswordPolicytTreeNodes(policy *CognitoPasswordRequirements) []*tree.Tree {
	var nodes []*tree.Tree

	if policy == nil {
		return nodes
	}

	root := tree.New().Root(styles.Primary.Render("Password Policy"))
	root.Child(styles.Primary.Render("Min Length:         ") + styles.Secondary.Render(fmt.Sprintf("%d", policy.MinLength)))
	root.Child(styles.Primary.Render("Require Numbers:    ") + styles.Secondary.Render(fmt.Sprintf("%t", policy.RequireNumbers)))
	root.Child(styles.Primary.Render("Require Symbols:    ") + styles.Secondary.Render(fmt.Sprintf("%t", policy.RequireSymbols)))
	root.Child(styles.Primary.Render("Require Uppercase:  ") + styles.Secondary.Render(fmt.Sprintf("%t", policy.RequireUppercase)))
	root.Child(styles.Primary.Render("Require Lowercase:  ") + styles.Secondary.Render(fmt.Sprintf("%t", policy.RequireLowercase)))

	nodes = append(nodes, root)
	return nodes
}

func (cognito CognitoSettings) JoinedSignInIdentifiers() string {
	if cognito.SignInIdentifiers != nil && len(*cognito.SignInIdentifiers) > 0 {
		return strings.Join(*cognito.SignInIdentifiers, ", ")
	}
	return "[not set]"
}

func (cognito CognitoSettings) JoinedSignUpAttributes() string {
	if cognito.SignUpAttributes != nil && len(*cognito.SignUpAttributes) > 0 {
		return strings.Join(*cognito.SignUpAttributes, ", ")
	}
	return "[not set]"
}

func (client CognitoAppClient) JoinedReturnUrls() string {
	if client.ReturnUrls != nil && len(*client.ReturnUrls) > 0 {
		return strings.Join(*client.ReturnUrls, ", ")
	}
	return "[not set]"
}
