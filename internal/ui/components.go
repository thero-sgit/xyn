package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

func statusCmdComponent() string {
	cmdLabel := lipgloss.NewStyle().
		Background(lipgloss.Color("#f0efef")).
		Foreground(lipgloss.Color("#121111")).
		Render(" /status ")

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

func createTextArea(width int) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = ` Try "how does <filename> work?"`
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	
	ta.SetHeight(2)
	ta.SetWidth(width - 11)
	ta.Focus()

	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle()
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#EEEEEE"))
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)

	return ta
}

func dirAndSessionLabel(directoryPath, sessionName string) string {
	directoryPath = defaultBg.
		PaddingLeft(2).
		Foreground(lipgloss.Color("#E47753")).
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
	label = lipgloss.NewStyle().
		PaddingLeft(2).
		Bold(true).
		Foreground(lipgloss.Color("241")).
		Render(label)

	value = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(2).
		Render(value)

	joined := lipgloss.JoinHorizontal(lipgloss.Left, label, value)

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(joined)
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

	filled = lipgloss.NewStyle().Foreground(lipgloss.Color("#d95b5b")).Render(filled)
	unFilled = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(unFilled)

	return lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		AlignHorizontal(lipgloss.Center).
		Render(filled + unFilled + fmt.Sprintf(" [%f%s] ", usage, "%"))
}

func prettyError(e error, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		PaddingRight(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				"└ " + lipgloss.NewStyle().Foreground(lipgloss.Color("#d94444")).Render("Oops! ") + " " + e.Error() + "; ",
				"Please check your internet connection",
			),
		)
}