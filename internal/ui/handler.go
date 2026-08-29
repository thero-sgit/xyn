package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
	"github.com/thero-sgit/xyn/internal/ai"
)

var CHandler Handler

var responseChan     = make(chan openai.ChatCompletionMessage)
var sessionInfoChan  = make(chan string)

type Handler struct {
	session      *ai.Session
}

func (h *Handler) handlePrompt(prompt string, m Model) (Model, tea.Cmd, tea.Cmd) {
	var newSessionCmd tea.Cmd

	if len(h.session.History) < 1 {
		newSessionCmd = awaitSessionInfo
		h.session.NameSession(prompt, &sessionInfoChan)
	}

	h.session.NewPrompt(prompt, &responseChan)

	m.agentActivity = newAgentBackgroundActivity("Working")

	m.prettyHistory = append(m.prettyHistory, newUserPrompt(prompt, m.viewport.Width))
	m.prettyHistory = append(m.prettyHistory, m.agentActivity.prettyString)
    m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))
    m.textarea.Blur()
	m.textarea.Reset()

    if m.isChatClear {
        if len(h.session.History) > 0 {
            m.isChatClear = false
        }
    }

	m.agentActivity.index = len(m.prettyHistory)-1
    m.isAgentWorking = true

	return m, awaitResponse, newSessionCmd
}

type sessionInfo struct {
	Data string
}

func awaitSessionInfo() tea.Msg {
	return func() tea.Msg {
		completion := <- sessionInfoChan
		return sessionInfo { Data:  completion }
	}()	
}

type response struct {
	Data string
}

func awaitResponse() tea.Msg {
	return func() tea.Msg {
		completion := <- responseChan
		return response { Data:  completion.Content }
	}()	
}