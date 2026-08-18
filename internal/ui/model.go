package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	width  		int
	height 		int
	textarea 	textarea.Model
	chatHistory []string
	isAgentWorking bool
	agentActivity  agentBackgroundActivityLabel
	isChatClear    bool
}

func (m model) sendPrompt(prompt string) model {
	prompt = newUserPrompt(prompt)

	m.chatHistory = append(m.chatHistory, prompt)
	m.chatHistory = append(m.chatHistory, m.agentActivity.prettyString)

	m.textarea.Blur()
	m.textarea.Reset()

	if m.isChatClear {
		if len(m.chatHistory) > 0 {
			m.isChatClear = false
		}
	}

	m.isAgentWorking = true

	return m
}

func InitialModel() model {
	return model{
		agentActivity: newAgentBackgroundActivity("working"),
		isAgentWorking: false,
		isChatClear: true,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea = createTextArea(msg.Width)

	case tickMsg:	
		if m.isAgentWorking {
			m.agentActivity = m.agentActivity.animate()
			m.chatHistory[len(m.chatHistory)-1] = m.agentActivity.prettyString
			cmd = doTick()
		}

		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "i":
			if !m.textarea.Focused() {
				m.textarea.Focus()
				m.textarea.Reset()

				return m, cmd
			}

		case "d":
			if m.isAgentWorking {
				m.chatHistory[len(m.chatHistory)-1] = subtleStyle.Render("Done.")
				m.isAgentWorking = false
			}			

			return m, nil

		case "enter":
			prompt := strings.Trim(m.textarea.Value(), " ")

			if m.isAgentWorking {
				m.isAgentWorking = false
			}

			if prompt != "" {
				m = m.sendPrompt(prompt)

				if !m.isChatClear && m.isAgentWorking {
					cmd = doTick()
				}
			}

			return m, cmd
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing screen..."
	}

	// Vertical height allocations
	headerHeight := 3
	footerHeight := 3
	middleHeight := m.height - headerHeight - footerHeight

	// Assemble rows
	middle := lipgloss.JoinHorizontal(lipgloss.Top, m.chatUi(middleHeight))
	
	return lipgloss.JoinVertical(lipgloss.Left, m.headerBar(headerHeight), middle, m.footerBar(footerHeight))
}