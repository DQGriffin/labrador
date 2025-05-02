package styles

import (
	"github.com/DQGriffin/labrador/internal/cli/types"
	"github.com/charmbracelet/lipgloss"
)

var primaryColors = lipgloss.AdaptiveColor{Light: "238", Dark: "253"}
var Primary = lipgloss.NewStyle().Foreground(primaryColors)

var secondaryColors = lipgloss.AdaptiveColor{Light: "230", Dark: "248"}
var Secondary = lipgloss.NewStyle().Foreground(secondaryColors)

var tertiaryColors = lipgloss.AdaptiveColor{Light: "232", Dark: "246"}
var Tertiary = lipgloss.NewStyle().Foreground(tertiaryColors).Italic(true)

var errorForeground = lipgloss.AdaptiveColor{Light: "232", Dark: "15"}
var errorBackground = lipgloss.AdaptiveColor{Light: "232", Dark: "88"}
var Error = lipgloss.NewStyle().Foreground(errorForeground).Background(errorBackground).Bold(true)

var warnForeground = lipgloss.AdaptiveColor{Light: "232", Dark: "11"}
var Warn = lipgloss.NewStyle().Foreground(warnForeground)

var headingForeground = lipgloss.AdaptiveColor{Light: "0", Dark: "0"}
var headingBackground = lipgloss.AdaptiveColor{Light: "232", Dark: "248"}
var Heading = lipgloss.NewStyle().Foreground(headingForeground).Background(headingBackground).Bold(true)

var PrimaryStyle = types.ConsoleStyle{
	Bold:               false,
	Italic:             false,
	Underline:          false,
	Margin:             0,
	Padding:            0,
	ForegroundColorHex: "#f5f5f5",
}

var WarnStyle = types.ConsoleStyle{
	Bold:               false,
	Italic:             false,
	Underline:          false,
	Margin:             0,
	Padding:            0,
	ForegroundColorHex: "#ffff00",
}

var ErrorStyle = types.ConsoleStyle{
	Bold:               true,
	Italic:             false,
	Underline:          false,
	Margin:             0,
	Padding:            0,
	ForegroundColorHex: "#ffffff",
	BackgroundColorHex: "#870000",
}

var QuietStyle = types.ConsoleStyle{
	Bold:               false,
	Italic:             false,
	Underline:          false,
	Margin:             0,
	Padding:            0,
	ForegroundColorHex: "#949494",
}

var ResourceHeadingStyle = types.ConsoleStyle{
	Bold:               false,
	Italic:             true,
	Underline:          false,
	Margin:             0,
	Padding:            1,
	ForegroundColorHex: "#949494",
}
