package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	width  		   int
	height 		   int
	textarea 	   textarea.Model
	viewport       viewport.Model
	chatHistory    []string
	isAgentWorking bool
	agentActivity  agentBackgroundActivityLabel
	isChatClear    bool
}

func (m model) sendPrompt(prompt string) model {
    prompt = newUserPrompt(prompt, m.viewport.Width)

    m.chatHistory = append(m.chatHistory, prompt)
    m.chatHistory = append(m.chatHistory, m.agentActivity.prettyString)

    m.viewport.SetContent(strings.Join(m.chatHistory, "\n"))

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
		agentActivity: newAgentBackgroundActivity("Working"),
		isAgentWorking: false,
		isChatClear: true,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		vpCmd tea.Cmd
		taCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea = createTextArea(msg.Width)

		// Calculate layout vertical space
        headerHeight := 3
        footerHeight := 3
        middleHeight := m.height - headerHeight - footerHeight
        promptBoxHeight := 4 // Account for textarea + padding

        vpWidth := m.width - 4
        vpHeight := middleHeight - promptBoxHeight

        if vpHeight < 1 {
            vpHeight = 1
        }

        // Initialize or update viewport dimensions
        m.viewport = viewport.New(vpWidth, vpHeight)
        m.viewport.SetContent(strings.Join(m.chatHistory, "\n"))

	case tickMsg:	
		if m.isAgentWorking {
			m.agentActivity.animate()
			m.chatHistory[len(m.chatHistory)-1] = m.agentActivity.prettyString
			m.viewport.SetContent(strings.Join(m.chatHistory, "\n"))
			m.viewport.GotoBottom()

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

		case "ctrl+d":
			if m.isAgentWorking {
				m.chatHistory[len(m.chatHistory)-1] = subtleStyle.Render("Done.")
				m.viewport.SetContent(strings.Join(m.chatHistory, "\n"))
				m.viewport.GotoBottom()
				m.isAgentWorking = false
			}			

			return m, nil

		case "alt+enter":
			if m.isAgentWorking {
				return m, cmd
			}

			prompt := strings.Trim(m.textarea.Value(), " ")		

			if prompt != "" {
				m = m.sendPrompt(prompt)
				m.viewport.GotoBottom()

				if !m.isChatClear && m.isAgentWorking {
					cmd = doTick()
				}
			}

			return m, cmd

		case "pgup", "pgdown", "up", "down":
			if !m.textarea.Focused() {
				m.viewport, vpCmd = m.viewport.Update(msg)
				return m, vpCmd	
			}
		}
	}

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(taCmd, vpCmd, cmd)
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