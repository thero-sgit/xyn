package ui

import (
	"strings"

	"github.com/thero-sgit/xyn/internal/ai"
)

var CHandler Handler

type Handler struct {
	session *ai.Session
}

func (h *Handler) handlePrompt(prompt string, m Model) Model {
	h.session.NewPrompt(prompt)

	vpContent := viewportContent(h.session.History, m.viewport.Width)
	vpContent  = append(vpContent, m.agentActivity.prettyString)
    m.viewport.SetContent(strings.Join(vpContent, "\n"))

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