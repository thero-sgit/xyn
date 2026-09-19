package ui

import (
    "context"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/thero-sgit/xyn/internal/ai"
)

var CHandler Handler

type Handler struct {
    session *ai.Session
}

func InitHandler() {
    CHandler = Handler{
        session: &ai.CSession,
    }
}

type errMsg struct{ err error }

func errorListener(ctx context.Context, errChan <-chan error) tea.Cmd {
    return func() tea.Msg {
        select {
        case e, ok := <-errChan:
            if !ok || e == nil {
                return nil
            }
            return errMsg{err: e}
        case <-ctx.Done():
            return errorListener(ctx, errChan)()
        }
    }
}

type response struct {
    chunk ai.Chunk
    c     chan ai.Chunk
}

func awaitResponse(ctx context.Context, c chan ai.Chunk) tea.Cmd {
    return func() tea.Msg {
        select {
        case completion, ok := <-c:
            if !ok {
                return nil
            }
            
            return response{chunk: completion, c: c}
        case <-ctx.Done():
            return nil
        }
    }
}

type sessionInfo struct {
    chunk ai.Chunk
    c     chan ai.Chunk
}

func awaitSessionInfo(ctx context.Context, c chan ai.Chunk) tea.Cmd {
    return func() tea.Msg {
        select {
        case chunk, ok := <-c:
            if !ok {
                return nil
            }
            return sessionInfo{chunk: chunk, c: c}
        case <-ctx.Done():
            return nil
        }
    }
}

func (h *Handler) handlePrompt(m Model) (Model, tea.Cmd) {
    responseChan := make(chan ai.Chunk)
    sessionInfoChan := make(chan ai.Chunk)
    errChan := make(chan error, 2)

    prompt := m.chatCentre.currentUserPrompt.message
    reqCtx := h.session.Ctx

    var cmds []tea.Cmd

    if len(h.session.History) < 1 {
        h.session.NameSession(prompt, sessionInfoChan, errChan)
        cmds = append(cmds, awaitSessionInfo(reqCtx, sessionInfoChan))
    }

    h.session.NewPrompt(prompt, responseChan, errChan)

    m.textarea.Blur()
    m.textarea.Reset()

    if m.chatCentre.isChatClear && len(h.session.History) > 0 {
        m.chatCentre.isChatClear = false
    }

    m.chatCentre.sendingPrompt = true

    cmds = append(cmds, errorListener(reqCtx, errChan), awaitResponse(reqCtx, responseChan))

    return m, tea.Batch(cmds...)
}