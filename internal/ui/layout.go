package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thero-sgit/xyn/internal/config"
)

func (m Model) headerBar(headerHeight int) string {
	version := subtleStyle.Render("v0.4.2")

	hemiStyle := lipgloss.NewStyle().
		Width(m.width/2).
		PaddingRight(2)

	leftHemi := hemiStyle.
		AlignHorizontal(lipgloss.Left).
		Render()

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

	inputBoxContent := lipgloss.JoinHorizontal(lipgloss.Top, prompt, m.textarea.View())

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width-6)

	leftHemiContent := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left, dirAndSessionLabel(config.Config.SanitizedWd, m.sessionName)),
		)

	modelLabel := labelValueBand("❋", config.Config.Model)
	inpLabels := lipgloss.NewStyle().PaddingLeft(1).MarginTop(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Left, modelLabel),
	)
	
	promptBox := boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, inputBoxContent, inpLabels),
	)

	activity := lipgloss.NewStyle().
		Width(m.width-4).
		Height(middleHeight-4).
		MarginLeft(2).
		MarginRight(2).
		Render(m.viewport.View())

	main := lipgloss.JoinVertical(lipgloss.Top, activity, leftHemiContent, promptBox)	

	return lipgloss.NewStyle().
		Width(m.width).
		Height(middleHeight).
		PaddingLeft(2).
		PaddingRight(2).
		Render(main)
}

func (m Model) footerBar(footerHeight int) string {
	joined := lipgloss.JoinHorizontal(lipgloss.Left,)

	return lipgloss.NewStyle().
		PaddingRight(2).
		PaddingLeft(2).
		Width(m.width).
		Height(footerHeight - 2).
		Render(joined)
}