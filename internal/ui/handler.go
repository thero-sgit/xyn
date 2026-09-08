package ui

import (
	"time"

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

func (h *Handler) handlePrompt(m Model) (Model, tea.Cmd) {	
	var newSessionCmd tea.Cmd

	prompt := m.chatCentre.currentUserPrompt.message

	if len(h.session.History) < 1 {
		newSessionCmd = awaitSessionInfo
		h.session.NameSession(prompt, sessionInfoChan, errChan)
	}

	h.session.NewPrompt(prompt, responseChan, errChan)

    m.textarea.Blur()
	m.textarea.Reset()

    if m.chatCentre.isChatClear {
        if len(h.session.History) > 0 {
            m.chatCentre.isChatClear = false
        }
    }

	m.chatCentre.sendingPrompt  = true

	return m, tea.Batch(awaitResponse, newSessionCmd, errorListener())
}

type errMsg error

func errorListener() tea.Cmd {
	return func() tea.Msg {
		e := <-errChan
		CHandler.session.SetNewContext()

		time.Sleep(3*time.Second)
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