package types

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/cli/styles"
	"github.com/DQGriffin/labrador/pkg/helpers"
	"github.com/charmbracelet/lipgloss/tree"
)

func (config LambdaData) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	for _, fn := range config.Functions {
		childNodes := fn.ToTreeNodes(verbose)
		nodes = append(nodes, childNodes...)
	}

	return nodes
}

func (lambda LambdaConfig) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(lambda.Name))

	if verbose {
		node.Child(styles.Primary.Render("Region:     ") + styles.Secondary.Render(*lambda.Region))
		node.Child(styles.Primary.Render("Code:       ") + styles.Secondary.Render(*lambda.Code))
		node.Child(styles.Primary.Render("Handler:    ") + styles.Secondary.Render(*lambda.Handler))
		node.Child(styles.Primary.Render("Runtime:    ") + styles.Secondary.Render(*lambda.Runtime))
		node.Child(styles.Primary.Render("Role ARN:   ") + styles.Secondary.Render(*lambda.RoleArn))
		node.Child(styles.Primary.Render("Memory:     ") + styles.Secondary.Render(fmt.Sprintf("%dmb", *lambda.MemorySize)))
		node.Child(styles.Primary.Render("Timeout:    ") + styles.Secondary.Render(fmt.Sprintf("%ds", *lambda.Timeout)))
		node.Child(styles.Primary.Render("On Delete:  ") + styles.Secondary.Render(helpers.PtrOrDefault(lambda.OnDelete, "delete")))
	}

	nodes = append(nodes, node)
	return nodes
}
