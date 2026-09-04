package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func newUserPrompt(prompt string, width int) string {
	return lipgloss.NewStyle().
		MarginBottom(1).
		Width(width-6).
		PaddingRight(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Background(accentColor).Render(" you $ "),
				lipgloss.NewStyle().PaddingLeft(1).Render(prompt),
			),
		)
}

func agentResponse(message string, width int) string {
	return lipgloss.NewStyle().
		MarginBottom(1).
		Width(width).
		Padding(1).
		Render(message)
}

type agentBackgroundActivityLabel struct {
	index          int
	prefixLabel    string
	rawString      string
	loaderIndex    int
	loaderFrames   []string
	loader         string
	prettyString   string
	doneStyle 	   lipgloss.Style
	intrptStyle    lipgloss.Style
	loaderStyle    lipgloss.Style
}

func newAgentBackgroundActivity(activity string) agentBackgroundActivityLabel {
	prefixLabel := bwLabelStyle.Render(" xyn $ ")

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	loader := frames[0]


	doneStyle   := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Italic(true)
	loaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E47753")).Bold(true)
	intrptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Bold(true)

	prettyString := fmt.Sprintf(
		"%s %s %s...",
		prefixLabel,
		loaderStyle.Render(loader),
		activity,
	)

	return agentBackgroundActivityLabel {
		prefixLabel:  prefixLabel,
		rawString:    activity,
		loaderIndex:  0,
		loaderFrames: frames,
		loader:       loader,
		prettyString: prettyString,
		doneStyle:    doneStyle,
		intrptStyle: intrptStyle,
		loaderStyle:  loaderStyle,
	}
}

func (abal *agentBackgroundActivityLabel) animate() {
	abal.loaderIndex = (abal.loaderIndex + 1) % len(abal.loaderFrames)

	abal.prettyString = fmt.Sprintf(
		"%s %s %s...",
		abal.prefixLabel,
		abal.loaderStyle.Render(abal.loaderFrames[abal.loaderIndex]),
		abal.rawString,
	)
}

func (abal *agentBackgroundActivityLabel) done(onErr bool) {
	var concl string
	if onErr{
		concl = abal.intrptStyle.Render("*")
	} else {
		concl = abal.doneStyle.Render("Thought process") + ">"
	}

	abal.prettyString = fmt.Sprintf(
		"%s %s",
		abal.prefixLabel,
		concl,
	)
}