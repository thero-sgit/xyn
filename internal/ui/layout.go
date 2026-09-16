package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/thero-sgit/xyn/internal/config"
	zone "github.com/lrstanley/bubblezone"
)

func (m Model) headerBar(headerHeight int) string {
	version := subtleStyle.Render("v0.4.2")

	return lipgloss.NewStyle().
		Width(m.width).
		Height(headerHeight).
		AlignVertical(lipgloss.Bottom).
		AlignHorizontal(lipgloss.Right).
		PaddingRight(2).
		Render(version)
}


func (m Model) chatUi(middleHeight int) string {
	prompt := lipgloss.NewStyle().
		Height(2).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render("  >  ")

	inputBoxContent := zone.Mark("text-area", lipgloss.JoinHorizontal(lipgloss.Top, prompt, m.textarea.View()))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width-6)

	leftHemiContent := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left, dirAndSessionLabel(config.Config.SanitizedWd, m.chatCentre.sessionName)),
		)

	addContext := lipgloss.NewStyle().
		MarginLeft(1).
		PaddingLeft(1).
		PaddingRight(1).
		Render("\uFF0B")

	addContextbutton := zone.Mark("add-context-btn", addContext)
	modelLabelButton := zone.Mark("model-selector-btn", labelValueBand("⬡", config.Config.Model))
	modeLabelButton  := zone.Mark("mode-toggle-btn", labelValueBand(config.Config.Mode.Icon, config.Config.Mode.Name))
	sendPromptButton := zone.Mark("send-prompt-btn", sendPromptState)

	btnContHalfStyle := lipgloss.NewStyle().Width((m.width-8)/2)

	inpLabels := lipgloss.NewStyle().PaddingLeft(1).MarginTop(1).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left, 
			btnContHalfStyle.Render(addContextbutton, modelLabelButton + " ", modeLabelButton + " ",),
			btnContHalfStyle.AlignHorizontal(lipgloss.Right).Render(sendPromptButton),
		),
	)
	
	promptBox := boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, inputBoxContent, inpLabels),
	)

	activity := lipgloss.NewStyle().
		Width(m.width-4).
		Height(middleHeight-8).
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
	joined := subtleStyle.Render("ctrl+e to write prompt   ? for help   esc/ctrl+c to quit")

	return lipgloss.NewStyle().
		PaddingRight(2).
		PaddingLeft(2).
		Width(m.width).
		AlignHorizontal(lipgloss.Right).
		Height(footerHeight).
		Render(joined)
}