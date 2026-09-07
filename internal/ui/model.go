package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type tickMsg50 time.Time
func doTick50() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg50(t)
	})
}

type tickMsg20 time.Time
func doTick20() tea.Cmd {
	return tea.Tick(20*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg20(t)
	})
}

type buttonTracker struct {
	modelButtonClicked bool
}

type Model struct {
	width  		   int
	height 		   int
	textarea 	   textarea.Model
	viewport       viewport.Model

	//
	chatCentre     chatCentre

	//
	slashCmdCon    slashCmdCtrl
	buttonTracker  buttonTracker
	
	//
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
		chatCentre: newChatCentre(),
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
		m.textarea = createTextArea(msg.Width, m.textarea.Value())

        headerHeight := 2
        footerHeight := 1
        middleHeight := m.height - headerHeight - footerHeight
        promptBoxHeight := 7

		m.viewport.Height = max(middleHeight - promptBoxHeight, 1)
		m.viewport.Width  = m.width - 4

		var vpContent string
		if m.Error != nil {
			vpContent = prettyError(m.Error, m.viewport.Width)
		} else if len(m.chatCentre.prettyHistory()) < 1 {
			vpContent = statusCmdComponent()
		} else {
			vpContent = strings.Join(m.chatCentre.prettyHistory(), "\n")
		}
        m.viewport.SetContent(vpContent)
		m.chatCentre.updatedWidths(m.viewport.Width)

	case tea.MouseMsg:
        if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			m.buttonTracker = buttonTracker{}
            z := zone.Get("model-selector")
            
            if z.InBounds(msg) {				
				m.buttonTracker.modelButtonClicked = true
                return m, nil
            }
        }

	case errMsg:
		// m.Error = msg
		// m.chatCentre.isAgentWorking = false
		// m.chatCentre.agentActivity.done(true)
		// m.chatCentre.prettyHistory() = append(m.chatCentre.prettyHistory(), prettyError(m.Error, m.viewport.Width))
		// m.viewport.SetContent(strings.Join(m.chatCentre.prettyHistory(), "\n"))

		return m, nil

	case tickMsg50:
		if m.Error != nil {
			return m, nil
		}

		m.chatCentre.history[m.chatCentre.currentAgentRes.index] = m.chatCentre.currentAgentRes
		m.viewport.SetContent(strings.Join(m.chatCentre.prettyHistory(), "\n"))

		if m.chatCentre.isAgentWorking {
			m.chatCentre.currentAgentRes.agentBgActivity.animate()
			cmd = doTick50()
		}

		return m, cmd

	case tickMsg20:
		if m.Error != nil {
			return m, nil
		}

		m.chatCentre.history[m.chatCentre.currentUserPrompt.index] = m.chatCentre.currentUserPrompt
		m.viewport.SetContent(strings.Join(m.chatCentre.prettyHistory(), "\n"))

		if m.chatCentre.isAgentWorking {
			m.chatCentre.currentUserPrompt.animate()
			cmd = doTick20()
		}

		return m, cmd

	case response:
		m.chatCentre.currentAgentRes.responseBuffer += msg.Data

		if msg.EOS {
			m.chatCentre.isAgentWorking = false
			m.chatCentre.currentAgentRes.responseBuffer = ""
			return m, nil
		}

		m.chatCentre.history[m.chatCentre.currentAgentRes.index] = m.chatCentre.currentAgentRes
		m.viewport.SetContent(strings.Join(m.chatCentre.prettyHistory(), "\n"))
		m.viewport.GotoBottom()

		return m, awaitResponse

	case sessionInfo:
		m.chatCentre.sessionName += msg.Data

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
				m.textarea.Focus()
				m.slashCmdCon.acceptingCmd = false
			} else {
				m.textarea.Blur()
			}				

			return m, cmd

		case "ctrl+x":
			if m.chatCentre.isAgentWorking {
				m.chatCentre.currentAgentRes.agentBgActivity.done(false)
				m.viewport.GotoBottom()				
				m.chatCentre.isAgentWorking = false

				return m, doTick50()
			}			

		case "alt+enter":
			if m.chatCentre.isAgentWorking {
				return m, cmd
			}

			var r tea.Cmd
			var s tea.Cmd

			prompt := strings.Trim(strings.TrimSpace(m.textarea.Value()), "\n")

			if prompt != "" {
				m, r, s = CHandler.handlePrompt(prompt, m)
				m.Error = nil
				m.viewport.GotoBottom()

				cmd = doTick50()
			}

			return m, tea.Batch(cmd, r, s, awaitResponse, errorListener())

		case "pgup", "pgdown", "up", "down":
			if !m.textarea.Focused() {
				m.viewport, vpCmd = m.viewport.Update(msg)
			} else {
				m.textarea, taCmd = m.textarea.Update(msg)
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

	if m.buttonTracker.modelButtonClicked {
		m.viewport.SetContent("MODELL!!")
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
	headerHeight := 2
	footerHeight := 1
	middleHeight := m.height - headerHeight - footerHeight

	// Assemble rows
	middle := lipgloss.JoinHorizontal(lipgloss.Top, m.chatUi(middleHeight))
	
	screen := lipgloss.JoinVertical(lipgloss.Left, m.headerBar(headerHeight), middle, m.footerBar(footerHeight))

	return zone.Scan(screen)
}