package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type chatCentre struct {
	sessionName       string
	isAgentWorking    bool
	agentActivity     agentBackgroundActivityLabel
	isChatClear       bool
	prettyHistory     []string
	responseBuffer    string
	currentUserPrompt userPrompt
}

func newChatCentre() chatCentre {
	return chatCentre{
		agentActivity: newAgentBackgroundActivity("Working"),
		isAgentWorking: false,
		isChatClear: true,
	}
}

type userPrompt struct {
	index     int
	width     int
	message   string
	waveIndex int
	pretty    string
}

func newUserPrompt(prompt string, width int) userPrompt {
	width = width-6
	pretty := lipgloss.NewStyle().
		MarginBottom(1).
		Width(width).
		PaddingRight(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Background(accentColor).Render(" you $ "),
				lipgloss.NewStyle().PaddingLeft(1).Render(prompt),
			),
		)

	return userPrompt{
		width:     width,
		message:   prompt,
		waveIndex: 0,
		pretty:    pretty,	
	}
}

func (up *userPrompt) animate() {
	up.waveIndex++

	var newPretty strings.Builder
	for i, v := range up.message {
		phase := float64(up.waveIndex)*0.2 - float64(i)*0.4
		offset := math.Sin(phase)

		style := colorForOffset(offset)
		newPretty.WriteString(style.Render(string(v)))
	}

	up.pretty = lipgloss.NewStyle().
		MarginBottom(1).
		Width(up.width).
		PaddingRight(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Background(accentColor).Render(" you $ "),
				lipgloss.NewStyle().PaddingLeft(1).Render(newPretty.String()),
			),
		)
}

func colorForOffset(offset float64) lipgloss.Style {
	level := 235 + int((offset + 1)/2*23)
	return lipgloss.NewStyle().Foreground(
		lipgloss.Color(fmt.Sprintf("%d", level)),
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