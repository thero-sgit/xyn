package ui

import "github.com/charmbracelet/lipgloss"

func dirAndSessionLabel(directoryPath, sessionName string) string {
	directoryPath = defaultBg.
		PaddingLeft(2).
		Foreground(lipgloss.Color("#E47753")).
		Render("~" + directoryPath)
	
	separator := defaultBg.
		Foreground(lipgloss.Color("#5DCAA5")).
		Render(" • ")

	sessionName = defaultBg.
		PaddingRight(2).
		Render(sessionName)

	joined := lipgloss.JoinHorizontal(lipgloss.Left, directoryPath, separator, sessionName)

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		PaddingLeft(2).
		PaddingRight(2).
		Render(joined)
}

func labelValueBand(label, value string) string {
	label = defaultBg.
		PaddingLeft(2).
		Foreground(lipgloss.Color("241")).
		Render(label + "    ")

	value = defaultBg.
		PaddingRight(2).
		Render("    " + value)

	joined := lipgloss.JoinHorizontal(lipgloss.Left, label, value)

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(joined)
}

func contextProgessBar(width int, progressPerc float32) string {
	widthAdj := width - 2
	bars := int(float32(widthAdj) * progressPerc)

	var filled string 

	for i := 0; i < bars; i++ {
		filled += "—"
	}

	var unFilled string 

	for i := 0; i < widthAdj - bars; i++ {
		unFilled += "—"
	}

	filled = lipgloss.NewStyle().Foreground(lipgloss.Color("#d95b5b")).Render(filled)
	unFilled = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(unFilled)

	return lipgloss.NewStyle().
		Bold(true).
		Width(width).
		AlignHorizontal(lipgloss.Center).
		Render(filled + unFilled)
}