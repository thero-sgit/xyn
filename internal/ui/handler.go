package ui

import (
	"strings"

	"github.com/thero-sgit/xyn/internal/ai"
)

var CHandler Handler

type Handler struct {
	session *ai.Session
	model   Model
}

func (h *Handler) handlePrompt(prompt string, m Model) Model {
	h.session.NewPrompt(prompt)

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

    m.isAgentWorking = true

	return m
}