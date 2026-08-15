package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// State defines Mochi's current expression/activity
type State int

const (
	StateIdle State = iota
	StateThinking
	StateSuccess
	StateError
	StateSleep
)

// Messages for async ticks and task completions
type tickMsg time.Time
type taskFinishedMsg struct{}

type model struct {
	width  int
	height int
	state   State
	frame   int
	message string

	textarea textarea.Model
}

// Mochi's thinking animation frames (active reaching arms)
var thinkingFrames = []string{
	"[  ⬢ _ ⬡  ]💬", "[  ⬡ _ ⬢  ]💬",
}

func createTextArea(width int) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = ` Try "how does <filename> work?"`
	ta.ShowLineNumbers = false
	
	// Remove internal prompt entirely so it never scrolls away
	ta.Prompt = ""
	
	ta.SetHeight(2)
	ta.SetWidth(width - 7)
	ta.Focus()

	// Style input elements with unified background color
	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle()
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#EEEEEE"))
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)

	return ta
}

func InitialModel() model {

	return model{
		state:   StateIdle,
		frame:   0,
		message: "Mochi is ready to help! Press 't' to trigger a task.",
	}
}

// Command to drive frame animation for the thinking state
func tickCmd() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Simulated asynchronous worker task
func simulateTaskCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(3 * time.Second)
		return taskFinishedMsg{}
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

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "i":
			m.state = StateIdle
			m.message = "Standing by!"
			return m, nil

		case "t":
			m.state = StateThinking
			m.message = "Working on it..."
			m.frame = 0
			return m, tea.Batch(tickCmd(), simulateTaskCmd())

		case "s":
			m.state = StateSuccess
			m.message = "Task finished flawlessly!"
			return m, nil

		case "e":
			m.state = StateError
			m.message = "Uh oh! Something went wrong..."
			return m, nil

		case "z":
			m.state = StateSleep
			m.message = "Mochi is taking a quick nap."
			return m, nil
		}

	case tickMsg:
		if m.state == StateThinking {
			m.frame = (m.frame + 1) % len(thinkingFrames)
			return m, tickCmd()
		}

	case taskFinishedMsg:
		m.state = StateSuccess
		m.message = "Async task completed successfully!"
		return m, nil
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