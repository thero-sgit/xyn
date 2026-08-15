package ui

import (

	"github.com/charmbracelet/lipgloss"
)

type userPrompt struct {
	representation string
}

func newUserPrompt(prompt string) userPrompt {
	representation := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E47753")).
		Render(" $ ") + lipgloss.NewStyle().Render(prompt)

	return userPrompt{ representation }
}