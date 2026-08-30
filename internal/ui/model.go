package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thero-sgit/xyn/internal/ai"
)

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type Model struct {
	width  		   int
	height 		   int
	textarea 	   textarea.Model
	viewport       viewport.Model

	// chat
	sessionName    string
	isAgentWorking bool
	agentActivity  agentBackgroundActivityLabel
	isChatClear    bool
	prettyHistory  []string
}

func InitialModel() Model {
	model := Model{
		agentActivity: newAgentBackgroundActivity("Working"),
		isAgentWorking: false,
		isChatClear: true,
	}

	CHandler = Handler{
		session: &ai.CSession,
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

        headerHeight := 3
        footerHeight := 3
        middleHeight := m.height - headerHeight - footerHeight
        promptBoxHeight := 4

        vpWidth := m.width - 4
        vpHeight := max(middleHeight - promptBoxHeight, 1)

        m.viewport = viewport.New(vpWidth, vpHeight)

		var vpContent string
		if len(m.prettyHistory) < 1 {
			vpContent = statusCmdComponent()
		} else {
			vpContent = strings.Join(m.prettyHistory, "\n")
		}
        m.viewport.SetContent(vpContent)

	case tickMsg:
		m.prettyHistory[m.agentActivity.index] = m.agentActivity.prettyString
		m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))

		if m.isAgentWorking {
			m.agentActivity.animate()
			cmd = tea.Batch(doTick(), awaitResponse)
		}

		return m, cmd

	case response:
		m.prettyHistory = append(m.prettyHistory, msg.Data)
		m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))

		return m, nil

	case sessionInfo:
		if msg.Data.Error != nil {
			m.isAgentWorking = false
			m.agentActivity.done(true)
			m.prettyHistory = append(m.prettyHistory, prettyError(msg.Data.Error, m.viewport.Width))
			m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))
		} else {
			m.sessionName = msg.Data.Name
		}		

		return m, nil

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

		case "ctrl+x":
			if m.isAgentWorking {
				m.agentActivity.done(false)
				m.viewport.GotoBottom()				
				m.isAgentWorking = false

				return m, doTick()
			}			

		case "alt+enter":
			if m.isAgentWorking {
				return m, cmd
			}

			prompt := strings.Trim(strings.TrimSpace(m.textarea.Value()), "\n")
			var r tea.Cmd
			var s tea.Cmd

			if prompt != "" {
				m, r, s = CHandler.handlePrompt(prompt, m)
				m.viewport.GotoBottom()

				cmd = doTick()
			}

			return m, tea.Batch(cmd, r, s)

		case "pgup", "pgdown", "up", "down":
			if !m.textarea.Focused() {
				m.viewport, vpCmd = m.viewport.Update(msg)
				return m, vpCmd	
			}


		case "/":
			m.viewport.SetContent(statusCmdComponent())
		}
	}

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(taCmd, vpCmd, cmd)
}

func (m Model) View() string {
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