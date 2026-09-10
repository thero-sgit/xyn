package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	defaultBg    = lipgloss.NewStyle().Background(lipgloss.Color("#1A1A1A"))
	bwLabelStyle = lipgloss.NewStyle().Background(lipgloss.Color("#f0efef")).Foreground(lipgloss.Color("#121111"))
	accentColor  = lipgloss.Color("#ff6f47")

	slashCommands = []slashCommand {
		{ name: "/help  ", pretty: "/help  ", comp: helpCmdComponent},
		{ name: "/status", pretty: "/status", comp: statusCmdComponent},
	}

	sendPromptBtnStates = []string {
		lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(accentColor).
		Render("\u27A4"),

		lipgloss.NewStyle().
		PaddingRight(1).
		Render("\u25A0"),
	}

	sendPromptState = sendPromptBtnStates[0]
)

type slashCommand struct{
	name   string
	pretty string
	comp   func() string
}

type slashCmdCtrl struct {
	ta               textinput.Model
	cmdLen           int
	acceptingCmd     bool
	highlighted 	 slashCommand
	highlightedIndex int	 
}

func newSlashCmdCtrl() slashCmdCtrl {
	ta := textinput.New()
	ta.Focus()

	return slashCmdCtrl{
		ta:               ta,
		cmdLen:           len(slashCommands),
		acceptingCmd:     true,
		highlightedIndex: 0,
	}
}

func (s *slashCmdCtrl) view() string {
	if len(s.ta.Value()) > 0 && !strings.HasPrefix(s.ta.Value(), "/") {
		s.ta.SetValue("/" + s.ta.Value())
		s.ta.CursorEnd()
	}

	cmds := []slashCommand{}
	for _, command := range slashCommands {
		if strings.HasPrefix(command.name, s.ta.Value()) {
			cmds = append(cmds, command)
		}
	}
	s.cmdLen = len(cmds)
	if s.highlightedIndex >= s.cmdLen || s.highlightedIndex < 0 {
		s.highlightedIndex = max(s.cmdLen - 1, 0)
	}
	
	if s.cmdLen >= 1 {
		cmds[s.highlightedIndex].pretty = lipgloss.NewStyle().Foreground(accentColor).Render("> ") + cmds[s.highlightedIndex].name
		s.highlighted = cmds[s.highlightedIndex]
	}	

	toDisplay := func() []string {
		var r []string
		for _, c := range cmds {r = append(r, c.pretty)}
		return  r
	}()

	s.ta.Placeholder = s.highlighted.name

	return lipgloss.NewStyle().Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			"command: " + s.ta.View(),
			strings.Join(toDisplay, "\n"),
			subtleStyle.MarginTop(1).Render("press 'esc' to close"),
		),
	)
}

func helpCmdComponent() string {
	cmdLabel := bwLabelStyle.Render(" /help ")

	cmdLabel = lipgloss.JoinHorizontal(
		lipgloss.Left, cmdLabel, 
		" ",
		"Useful '/' Commands",
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		cmdLabel,
	)
}

func statusCmdComponent() string {
	cmdLabel := bwLabelStyle.Render(" /status ")

	tokenUsage := "Token usage (34,500 / 200,000 tokens — 17.2%)"
	header 	   := lipgloss.JoinHorizontal(lipgloss.Left, cmdLabel, " ", tokenUsage)

	tokenUsageProgressBar := tokenUsageProgessBar(17.2)		
		

	tokenUsageBreakdown := fmt.Sprint(
		"└ Breakdown by Category:\n",
		"	  Conversation History ......... 18,200 tokens (52.8%)\n",
		"	  Loaded Files (3) ............. 9,400 tokens (27.2%)\n",
		"	  System & Instructions ........ 4,100 tokens (11.9%)\n",
		"	  Tools Use & Skills ........... 2,800 tokens ( 8.1%)\n",
	)

	context := fmt.Sprint(
		"└ Context:\n",
		"	 • Model (openai/gpt-oss-120b) Token Limit Resets in 14h 2m 33s\n",
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		tokenUsageProgressBar,
		tokenUsageBreakdown,
		context,
	)
}

func createTextArea(width int, value string) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = ` Try "how does <filename> work?"`
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	
	ta.SetHeight(2)
	ta.SetWidth(width - 11)
	ta.SetValue(value)

	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle()
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#EEEEEE"))

	return ta
}

func dirAndSessionLabel(directoryPath, sessionName string) string {
	directoryPath = defaultBg.
		PaddingLeft(2).
		Foreground(lipgloss.Color("#9dbdb7")).
		Render("~" + directoryPath)
	
	sColor := func() string {
		if sessionName == "" {
			return "#b1b2b1"
		}

		return "#5DCAA5"
	}()

	separator := defaultBg.
		Foreground(lipgloss.Color(sColor)).
		Render(" • ")

	sessionName = defaultBg.
		PaddingRight(2).
		Render(sessionName)

	joined := lipgloss.JoinHorizontal(lipgloss.Left, directoryPath, separator, sessionName)

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		PaddingRight(2).
		Render(joined)
}

func labelValueBand(label, value string) string {
	label = defaultBg.
		PaddingLeft(1).
		Bold(true).
		Foreground(lipgloss.Color("#9dbdb7")).
		Render(label)

	value = defaultBg.
		PaddingLeft(1).
		PaddingRight(1).
		Render(value)

	joined := lipgloss.JoinHorizontal(lipgloss.Left, label, value)

	output := defaultBg.
		AlignHorizontal(lipgloss.Left).
		Render(joined)

	return output
}

func tokenUsageProgessBar(usage float32) string {
	widthAdj  := 25
	usageBars := int(usage * (float32(widthAdj)/100))

	var filled string 

	for i := 0; i < usageBars; i++ {
		filled += "■ "
	}

	var unFilled string 

	for i := 0; i < widthAdj - usageBars; i++ {
		unFilled += "□ "
	}

	filled = lipgloss.NewStyle().Foreground(accentColor).Render(filled)
	unFilled = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(unFilled)

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		AlignHorizontal(lipgloss.Center).
		Render(filled + unFilled + fmt.Sprintf(" [%f%s] ", usage, "%"))
}

func retryButton() string {
	b := lipgloss.NewStyle().
		Italic(true).
		Underline(true).
		Foreground(lipgloss.Color("#d94444")).
		Render("Retry "+ "\u21BB")

	return zone.Mark("retry-button", b)
}