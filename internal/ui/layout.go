package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	defaultBg    = lipgloss.NewStyle().Background(lipgloss.Color("#1A1A1A"))
)

func (m Model) headerBar(headerHeight int) string {
	version := subtleStyle.Render("v0.4.2")

	hemiStyle := lipgloss.NewStyle().
		Width(m.width/2).
		PaddingRight(2)

	leftHemiContent := lipgloss.JoinHorizontal(lipgloss.Left, dirAndSessionLabel(workingDirSanitized(), "refactor auth"))

	if m.width <= 95 {
		leftHemiContent = lipgloss.JoinHorizontal(lipgloss.Left,)
	}

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


func (m Model) chatUi(middleHeight int) string {
	prompt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#67AB9F")).
		Height(2).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render("  >  ")

	inputBoxContent := lipgloss.JoinHorizontal(lipgloss.Top, prompt, m.textarea.View(),)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width-4).
		MarginLeft(2).
		MarginRight(2)

	promptBox := boxStyle.Render(inputBoxContent)

	activity := lipgloss.NewStyle().
		Width(m.width -4).
		Height(middleHeight - 4).
		MarginLeft(2).
		MarginRight(2).
		Render(m.viewport.View())

	main := lipgloss.JoinVertical(lipgloss.Top, activity, promptBox)	

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Width(m.width).
		Height(middleHeight - 1).
		Render(main)
}

func (m Model) footerBar(footerHeight int) string {
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