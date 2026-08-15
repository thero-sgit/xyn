package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	boldStyle    = lipgloss.NewStyle().Bold(true)
	defaultBg    = lipgloss.NewStyle().Background(lipgloss.Color("#1A1A1A"))
)

func (m model) headerBar(headerHeight int) string {
	logo := boldStyle.Render("xyn") + boldStyle.Foreground(lipgloss.Color("#FF6666")).Render("_")
	version := subtleStyle.Render("v0.4.2")

	hemiStyle := lipgloss.NewStyle().
		Width(m.width/2).
		PaddingLeft(2).
		PaddingRight(2)

	leftHemiContent := lipgloss.JoinHorizontal(lipgloss.Left, logo, dirAndSessionLabel("/Desktop/projects", "refactor auth"))

	leftHemi := hemiStyle.
		AlignHorizontal(lipgloss.Left).
		
		Render(leftHemiContent)

	rightHemi := hemiStyle.
		AlignHorizontal(lipgloss.Right).
		Render(version)

	bar := lipgloss.JoinHorizontal(
		lipgloss.Top, 
		leftHemi,
		rightHemi,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		MarginTop(1).
		Height(headerHeight).
		AlignVertical(lipgloss.Center).
		Render(bar)
}


func (m model) chatUi(middleHeight int) string {
	prompt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B")).
		Background(lipgloss.Color("#1A1A1A")).
		Height(2).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render("  ⟩  ")

	inputBoxContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		prompt,
		m.textarea.View(),
	)

	boxStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1A1A1A")).
		Width(m.width-4).
		Padding(1).
		MarginLeft(2).
		MarginRight(2)

	promptBox := boxStyle.Render(inputBoxContent)

	activity := lipgloss.NewStyle().
		Width(m.width).
		Height(middleHeight - 4).
		Render()

	main := lipgloss.JoinVertical(lipgloss.Top, activity, promptBox)

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Width(m.width).
		Height(middleHeight - 1).
		Render(main)
}

func (m model) footerBar(footerHeight int) string {
	modelLabel := labelValueBand("model", "kimi-k2.5")
	sessionLabel := labelValueBand("session", "14m · 200k")

	joined := lipgloss.JoinHorizontal(lipgloss.Left, modelLabel, " ", sessionLabel)

	return lipgloss.NewStyle().
		PaddingRight(2).
		PaddingLeft(2).
		Width(m.width).
		Height(footerHeight - 2).
		Render(joined)
}