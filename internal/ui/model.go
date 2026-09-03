package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
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
	responseBuffer string

	//
	slashCmdCon    slashCmdCtrl
	
	Error          error
}

func InitialModel() Model {
	vp := viewport.New(0, 0)

	vp.KeyMap = viewport.KeyMap{
		Up:   key.NewBinding(key.WithKeys("up")),                                                                                                                                                                       
		Down: key.NewBinding(key.WithKeys("down")),                                                                                                                                                                         
		PageUp:  key.NewBinding(key.WithKeys("pgup")),                                                                                                                                                                                 
		PageDown: key.NewBinding(key.WithKeys("pgdn")),                                                                                                                                                                                                
		HalfPageUp:    key.NewBinding(),                                                                                                                                                                                                  
		HalfPageDown:  key.NewBinding(),
	}


	model := Model{
		agentActivity: newAgentBackgroundActivity("Working"),
		isAgentWorking: false,
		isChatClear: true,
		viewport: vp,
	}

	InitHandler()

	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, errorListener())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		vpCmd tea.Cmd
		taCmd tea.Cmd
		ciCmd tea.Cmd
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

		m.viewport.Height = max(middleHeight - promptBoxHeight, 1)
		m.viewport.Width  = m.width - 4

		var vpContent string
		if m.Error != nil {
			vpContent = prettyError(m.Error, m.viewport.Width)
		} else if len(m.prettyHistory) < 1 {
			vpContent = statusCmdComponent()
		} else {
			vpContent = strings.Join(m.prettyHistory, "\n")
		}
        m.viewport.SetContent(vpContent)

	case errMsg:
		m.Error = msg
		m.isAgentWorking = false
		m.agentActivity.done(true)
		m.prettyHistory = append(m.prettyHistory, prettyError(m.Error, m.viewport.Width))
		m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))

		return m, nil

	case tickMsg:
		if m.Error != nil {
			return m, nil
		}

		m.prettyHistory[m.agentActivity.index] = m.agentActivity.prettyString
		m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))

		if m.isAgentWorking {
			m.agentActivity.animate()
			cmd = doTick()
		}

		return m, cmd

	case response:
		m.responseBuffer += msg.Data

		if msg.EOS {
			m.isAgentWorking = false
			m.responseBuffer = ""
			return m, nil
		}

		m.prettyHistory[len(m.prettyHistory)-1] = agentResponse(m.responseBuffer, m.viewport.Width)
		m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))
		m.viewport.GotoBottom()

		return m, awaitResponse

	case sessionInfo:
		m.sessionName += msg.Data

		if msg.EOS {
			return m, nil
		}

		return m, awaitSessionInfo

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+e":
			if !m.textarea.Focused() {
				m.textarea.Reset()
				m.textarea.Focus()

			} else {
				m.textarea.Blur()
			}				

			return m, cmd

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

			var r tea.Cmd
			var s tea.Cmd

			prompt := strings.Trim(strings.TrimSpace(m.textarea.Value()), "\n")

			if prompt != "" {
				m, r, s = CHandler.handlePrompt(prompt, m)
				m.Error = nil
				m.viewport.GotoBottom()

				cmd = doTick()
			}

			return m, tea.Batch(cmd, r, s, awaitResponse, errorListener())

		case "pgup", "pgdown", "up", "down":
			if !m.textarea.Focused() {
				m.viewport, vpCmd = m.viewport.Update(msg)
			}

			if m.slashCmdCon.acceptingCmd {
				if msg.String() == "up" {
					m.slashCmdCon.highlightedIndex = (m.slashCmdCon.highlightedIndex - 1) % m.slashCmdCon.cmdLen
				}

				if msg.String() == "down" {
					m.slashCmdCon.highlightedIndex = (m.slashCmdCon.highlightedIndex + 1) % m.slashCmdCon.cmdLen
				}

				if len(m.slashCmdCon.ta.Value()) <= 1 { m.slashCmdCon.ta.Reset() }
				m.viewport.SetContent(lipgloss.JoinVertical(
					lipgloss.Top,
					m.slashCmdCon.view(),
				))
			}

			return m, vpCmd

		case "/":
			if !m.textarea.Focused() {
				m.slashCmdCon.acceptingCmd = true
				m.slashCmdCon              = newSlashCmdCtrl()			
			}

		case "enter":
			if m.slashCmdCon.acceptingCmd && m.slashCmdCon.cmdLen > 0 && !m.textarea.Focused() {
				m.viewport.SetContent(m.slashCmdCon.highlighted.comp())
				m.slashCmdCon.acceptingCmd = false		
			}
		}
	}

	if m.slashCmdCon.acceptingCmd {
		m.viewport.SetContent(lipgloss.JoinVertical(
			lipgloss.Top,
			m.slashCmdCon.view(),
		))
	}

	m.textarea, taCmd       = m.textarea.Update(msg)
	m.viewport, vpCmd       = m.viewport.Update(msg)
	m.slashCmdCon.ta, ciCmd = m.slashCmdCon.ta.Update(msg)
	return m, tea.Batch(taCmd, vpCmd, ciCmd, cmd)
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