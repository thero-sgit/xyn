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