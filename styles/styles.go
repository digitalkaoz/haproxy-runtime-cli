package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var AppName = "haproxy-runtime-cli"

var ActiveColor = lipgloss.CompleteAdaptiveColor{
	Light: lipgloss.CompleteColor{TrueColor: "#4c00c9", ANSI256: "93", ANSI: "5"},
	Dark:  lipgloss.CompleteColor{TrueColor: "#4c00c9", ANSI256: "93", ANSI: "5"},
}

var SecondaryColor = lipgloss.CompleteAdaptiveColor{
	Light: lipgloss.CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
	Dark:  lipgloss.CompleteColor{TrueColor: "#ffffff", ANSI256: "255", ANSI: "15"},
}

var ErrorColor = lipgloss.CompleteAdaptiveColor{
	Light: lipgloss.CompleteColor{TrueColor: "#ff0000", ANSI256: "160", ANSI: "9"},
	Dark:  lipgloss.CompleteColor{TrueColor: "#ff0000", ANSI256: "160", ANSI: "9"},
}
var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Background(ActiveColor).
	Foreground(SecondaryColor).
	Padding(1, 1, 1, 1).
	Underline(true)

var ActiveStyle = lipgloss.NewStyle().
	Foreground(ActiveColor).Bold(true)

var PageStyle = lipgloss.NewStyle().Margin(1, 2, 0, 2)

var ComplementStyle = lipgloss.NewStyle().
	Foreground(SecondaryColor).Bold(false)

var ResponseStyle = lipgloss.NewStyle().
	Foreground(SecondaryColor).Border(lipgloss.NormalBorder(), true, false, false, false).PaddingTop(1)

var ErrorStyle = lipgloss.NewStyle().
	Foreground(ErrorColor).Bold(true).Underline(true).Padding(1)
