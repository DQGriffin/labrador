package types

import "github.com/charmbracelet/lipgloss/tree"

func (stage Stage) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree

	if stage.Type == "s3" {
		for _, s3Config := range stage.Buckets {
			childNodes := s3Config.ToTreeNodes(verbose)
			nodes = append(nodes, childNodes...)
		}
	} else if stage.Type == "iam-role" {
		for _, config := range stage.IamRoles {
			childNodes := config.ToTreeNodes(verbose)
			nodes = append(nodes, childNodes...)
		}
	} else if stage.Type == "lambda" {
		for _, config := range stage.Functions {
			childNodes := config.ToTreeNodes(verbose)
			nodes = append(nodes, childNodes...)
		}
	} else if stage.Type == "api" {
		for _, config := range stage.Gateways {
			childNodes := config.ToTreeNodes(verbose)
			nodes = append(nodes, childNodes...)
		}
	}

	return nodes
}
