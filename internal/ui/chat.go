package ui

import (
	"os/user"

	"github.com/charmbracelet/lipgloss"
)

var currentUser, currentUserRrr  = user.Current()
var username = func() string {
		if currentUserRrr != nil {
			return "user"
		}
		return currentUser.Username
}()

func newUserPrompt(prompt string) string {
	return lipgloss.NewStyle().
		MarginBottom(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				subtleStyle.Render(username) + lipgloss.NewStyle().Foreground(lipgloss.Color("#E47753")).Render(" $ "),
				lipgloss.NewStyle().Render(prompt),
			),
		)
}

type agentBackgroundActivityLabel struct {
	rawString    string
	colorIndex   int
	colors       []string
	prettyString string
	style        lipgloss.Style
}

func newAgentBackgroundActivity(activity string) agentBackgroundActivityLabel {
	colors := []string {
		"#ffc0ac",
		"#ffa184",
		"#ff7e57",
		"#ff6b3a",
	}

	italicStyle := lipgloss.NewStyle().Italic(true)

	prettyString := italicStyle.
		Foreground(lipgloss.Color(colors[0])).
		Render(activity)

	return agentBackgroundActivityLabel {
		rawString: activity,
		colorIndex: 0,
		colors: colors,
		prettyString: prettyString,
		style: italicStyle,
	}
}

func (abal agentBackgroundActivityLabel) animate() agentBackgroundActivityLabel {
	abal.colorIndex = (abal.colorIndex + 1) % len(abal.colors)

	abal.prettyString = abal.style.
		Foreground(lipgloss.Color(abal.colors[abal.colorIndex])).
		Render(abal.rawString)

	return abal
}