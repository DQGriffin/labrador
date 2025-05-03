package types

import (
	"github.com/charmbracelet/lipgloss/tree"
)

type TreeInspectable interface {
	ToTreeNodes(verbose bool) []*tree.Tree
}
