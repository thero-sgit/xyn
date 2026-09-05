package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thero-sgit/xyn/internal/ai"
)

var CHandler Handler

var responseChan     = make(chan ai.Chunk)
var sessionInfoChan  = make(chan ai.Chunk)
var errChan          = make(chan error) 

type Handler struct {
	session      *ai.Session
}

func InitHandler() {
	CHandler = Handler{
		session: &ai.CSession,
	}
}

func (h *Handler) handlePrompt(prompt string, m Model) (Model, tea.Cmd, tea.Cmd) {	
	var newSessionCmd tea.Cmd

	if len(h.session.History) < 1 {
		newSessionCmd = awaitSessionInfo
		h.session.NameSession(prompt, sessionInfoChan, errChan)
	}

	h.session.NewPrompt(prompt, responseChan, errChan)

	m.chatCentre.agentActivity = newAgentBackgroundActivity("Working")
	up := newUserPrompt(prompt, m.viewport.Width)

	m.chatCentre.prettyHistory = append(m.chatCentre.prettyHistory, up.pretty)
	m.chatCentre.prettyHistory = append(m.chatCentre.prettyHistory, m.chatCentre.agentActivity.prettyString)
	m.chatCentre.prettyHistory = append(m.chatCentre.prettyHistory, m.chatCentre.responseBuffer)

    m.viewport.SetContent(strings.Join(m.chatCentre.prettyHistory, "\n"))
    m.textarea.Blur()
	m.textarea.Reset()

    if m.chatCentre.isChatClear {
        if len(h.session.History) > 0 {
            m.chatCentre.isChatClear = false
        }
    }

	m.chatCentre.currentUserPrompt       = up
	m.chatCentre.currentUserPrompt.index = len(m.chatCentre.prettyHistory) - 3
	m.chatCentre.agentActivity.index     = len(m.chatCentre.prettyHistory) - 2
    m.chatCentre.isAgentWorking          = true

	return m, awaitResponse, newSessionCmd
}

type errMsg error

func errorListener() tea.Cmd {
	return func() tea.Msg {
		e := <-errChan

		CHandler.session.SetNewContext()
		return errMsg(e)
	}
}

type sessionInfo ai.Chunk

func awaitSessionInfo() tea.Msg {
	return func() tea.Msg {
		select {
		case <-CHandler.session.Ctx.Done():
			return nil

		case chunk, ok := <- sessionInfoChan:
			if !ok {
				return nil
			}
			return sessionInfo(chunk)
		}		
	}()	
}

type response ai.Chunk

func awaitResponse() tea.Msg {
	return func() tea.Msg {
		select {
		case <-CHandler.session.Ctx.Done():
			return nil

		case completion, ok := <-responseChan:
			if !ok {
				return nil
			}
			return response(completion)
		}
	}()
}