package types

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/cli/styles"
	"github.com/DQGriffin/labrador/pkg/helpers"

	"github.com/charmbracelet/lipgloss/tree"
)

func (config S3Config) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	for _, bucket := range config.Buckets {
		childNodes := bucket.ToTreeNodes(verbose)
		nodes = append(nodes, childNodes...)
	}

	return nodes
}

func (s3 S3Settings) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(*s3.Name))

	if verbose {
		node.Child(styles.Primary.Render("Region:               ") + styles.Secondary.Render(helpers.PtrOrDefault(s3.Region, "[region not set]")))
		node.Child(styles.Primary.Render("Versioning:           ") + styles.Secondary.Render(fmt.Sprintf("%t", helpers.PtrOrDefault(s3.Versioning, false))))
		node.Child(styles.Primary.Render("Block Public Access:  ") + styles.Secondary.Render(fmt.Sprintf("%t", helpers.PtrOrDefault(s3.BlockPublicAccess, true))))
		node.Child(styles.Primary.Render("On Delete:            ") + styles.Secondary.Render(helpers.PtrOrDefault(s3.OnDelete, "delete")))
	}

	nodes = append(nodes, node)
	return nodes
}
