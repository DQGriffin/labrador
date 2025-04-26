package styles

import "github.com/charmbracelet/lipgloss"

var primaryColors = lipgloss.AdaptiveColor{Light: "238", Dark: "253"}
var Primary = lipgloss.NewStyle().Foreground(primaryColors)

var secondaryColors = lipgloss.AdaptiveColor{Light: "230", Dark: "248"}
var Secondary = lipgloss.NewStyle().Foreground(secondaryColors)

var tertiaryColors = lipgloss.AdaptiveColor{Light: "232", Dark: "246"}
var tertiaryColorsB = lipgloss.AdaptiveColor{Light: "230", Dark: "248"}
var Tertiary = lipgloss.NewStyle().Foreground(tertiaryColors).Italic(true)
