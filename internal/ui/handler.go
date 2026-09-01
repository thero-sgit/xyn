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

	m.agentActivity = newAgentBackgroundActivity("Working")

	m.prettyHistory = append(m.prettyHistory, newUserPrompt(prompt, m.viewport.Width))
	m.prettyHistory = append(m.prettyHistory, m.agentActivity.prettyString)
	m.prettyHistory = append(m.prettyHistory, m.responseBuffer)
    m.viewport.SetContent(strings.Join(m.prettyHistory, "\n"))
    m.textarea.Blur()
	m.textarea.Reset()

    if m.isChatClear {
        if len(h.session.History) > 0 {
            m.isChatClear = false
        }
    }

	m.agentActivity.index = len(m.prettyHistory)-2
    m.isAgentWorking = true

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