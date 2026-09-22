package ui

import (
	"fmt"
	"log"
	"math"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/thero-sgit/xyn/internal/config"
)

// -- CHAT CENTRE CONTROL

type chatCentre struct {
	sessionName       string
	isAgentWorking    bool
	isChatClear       bool
	history           []chatItem
	currentUserPrompt userPrompt
	currentAgentRes   agentResponse
	sendingPrompt     bool
}

func newChatCentre() chatCentre {
	return chatCentre{
		isAgentWorking: false,
		isChatClear: true,
	}
}

func (c *chatCentre) prettyHistory() []string {
	var h []string
	for _, item := range c.history {
		h = append(h, item.getPretty())
	}

	return h
}

func (c *chatCentre) updatedWidths(width int) {
	var updatedHistory []chatItem
	for _, v := range c.history {
		updatedHistory = append(updatedHistory, v.updated(width))
	}

	c.history = updatedHistory
}

// -- CHAT ITEMS CONTROL
type chatItem interface {
	getPretty() string
	updated(width int) chatItem
}

type userPrompt struct {
	index       int
	width       int
	message     string
	waveIndex 	int
	pretty    	string
}

func (up userPrompt) getPretty() string {
	return up.pretty
}

func (up userPrompt) updated(width int) chatItem {
	updated := newUserPrompt(up.message, width)
	updated.index = up.index
	updated.sent()

	return updated
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

func (up *userPrompt) sent() {
	p := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Background(accentColor).Render(" you $ "),
		lipgloss.NewStyle().PaddingLeft(1).Render(up.message),
	)

	up.pretty = lipgloss.NewStyle().
		MarginBottom(1).
		Width(up.width).
		PaddingRight(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				p,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#02e696")).PaddingRight(2).Render("\u2713"),
			),
		)
}

func (up *userPrompt) err() {
	p := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Background(accentColor).Render(" you $ "),
		lipgloss.NewStyle().Strikethrough(true).PaddingLeft(1).Render(up.message),
	)

	ed := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().MarginRight(4).Render(
			subtleStyle.Render("failed to send"),
		),
		retryButton(),
	)

	up.pretty = lipgloss.NewStyle().
		MarginBottom(1).
		Width(up.width).
		PaddingRight(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				p,
				ed,
			),
		)
}

func colorForOffset(offset float64) lipgloss.Style {
	level := 235 + int((offset + 1)/2*23)
	return lipgloss.NewStyle().Foreground(
		lipgloss.Color(fmt.Sprintf("%d", level)),
	)
}


type agentResponse struct {
	index           int
	width           int
	agentBgActivity *agentBackgroundActivityLabel
	reasoningBuffer string
	responseBuffer  string
	pretty          string
	showReasoning   bool
}

func (a agentResponse) updated(width int) chatItem {
	updated := newAgentResponse(width)
	updated.responseBuffer = a.responseBuffer
	updated.agentBgActivity = a.agentBgActivity
	updated.reasoningBuffer = a.reasoningBuffer

	return updated
}

func (a agentResponse) getPretty() string {
	content, err := config.Config.MdRenderer.Render(a.responseBuffer)
	if err != nil {
		log.Fatal(err)
	}

	if a.showReasoning {
		reasoningContent := subtleStyle.Width(a.width-2).Render(a.reasoningBuffer)
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			reasoningContent,
			"\n",
			content,
		)
	}

	return lipgloss.NewStyle().Width(a.width).MarginBottom(1).Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			a.agentBgActivity.prettyString,
			content,
		),
	)
}

func newAgentResponse(width int) agentResponse {
	agentBgActivity := newAgentBackgroundActivity("Thinking")

	pretty := lipgloss.NewStyle().Width(width).MarginBottom(1).Render(lipgloss.JoinVertical(
			lipgloss.Top,
			agentBgActivity.prettyString,
			lipgloss.NewStyle().
			Render(""),
		),
	)

	return agentResponse {
		width: width,
		agentBgActivity: &agentBgActivity,
		pretty: pretty,
	}
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
		intrptStyle:  intrptStyle,
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
		concl = abal.doneStyle.Render("Thought process") + " >"
	}

	prettyString := fmt.Sprintf(
		"%s %s",
		abal.prefixLabel,
		concl,
	)

	abal.prettyString = zone.Mark("thought-process-btn", prettyString)
}