package ai

import (
	"context"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var CSession Session

type Session struct {
	Ctx    context.Context
	cancl  context.CancelFunc
	groq   *groq

	History []openai.ChatCompletionMessage
}

func (s *Session) SetNewContext() {
	ctx, cancel := context.WithCancel(context.Background())
	s.Ctx = ctx
	s.cancl = cancel
}

func (s *Session) NameSession(message string, c *chan string, errChan *chan error) {
	req := openai.ChatCompletionRequest{
		Model: "openai/gpt-oss-20b",
		Messages:  []openai.ChatCompletionMessage{
			{Role: "system", Content: "Generate a 3-5 word title for this prompt. Return ONLY the title."},
			{Role: "user", Content: message},
		},
	}

	go func() {	
		resp, err := s.groq.client.CreateChatCompletion(s.Ctx, req)
		if err != nil {
			s.handleError(err, errChan)
			return
		}

		*c <- strings.TrimSpace(resp.Choices[0].Message.Content)
	}()
}

func (s *Session) NewPrompt(message string, c *chan openai.ChatCompletionMessage, errChan *chan error) {
	s.History = append(s.History, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleUser,
			Content:    message,
	})

	go func() {
        resp, err := s.groq.handlePrompt(s.Ctx, s.History)
        if err != nil {
            s.handleError(err, errChan)
            return
        }
		
        s.History = append(s.History, resp)
        *c <- resp
    }()
}

func (s *Session) handleError(err error, errChan *chan error) {
	s.cancl()
	*errChan <- err
}

func InitSession() {
	CSession = Session{
		groq: newGroq(),
	}

	CSession.SetNewContext()
}