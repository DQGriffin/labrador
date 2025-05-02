package types

import "github.com/charmbracelet/lipgloss"

type ConsoleStyle struct {
	Bold               bool
	Italic             bool
	Underline          bool
	Align              string
	Margin             int
	Padding            int
	ForegroundColorHex string
	BackgroundColorHex string
}

func (cs ConsoleStyle) ToLipglossStyle() lipgloss.Style {
	s := lipgloss.NewStyle()

	if cs.Bold {
		s = s.Bold(true)
	}
	if cs.Italic {
		s = s.Italic(true)
	}
	if cs.Underline {
		s = s.Underline(true)
	}
	if cs.ForegroundColorHex != "" {
		s = s.Foreground(lipgloss.Color(cs.ForegroundColorHex))
	}
	if cs.BackgroundColorHex != "" {
		s = s.Background(lipgloss.Color(cs.BackgroundColorHex))
	}
	switch cs.Align {
	case "center":
		s = s.Align(lipgloss.Center)
	case "right":
		s = s.Align(lipgloss.Right)
	default:
		s = s.Align(lipgloss.Left)
	}
	if cs.Margin > 0 {
		s = s.Margin(cs.Margin)
	}
	if cs.Padding > 0 {
		s = s.Padding(cs.Padding)
	}

	return s
}
