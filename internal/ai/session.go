package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

var CSession Session

type Session struct {
	ctx  context.Context
	groq *groq
	name string

	History []openai.ChatCompletionMessage
}

func (s *Session) nameSession(message string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: "llama-3.1-8b-instant",
		Messages:  []openai.ChatCompletionMessage{
			{Role: "system", Content: "Generate a 3-5 word title for this prompt. Return ONLY the title."},
			{Role: "user", Content: message},
		},
		MaxTokens: 15,
	}

	resp, err := s.groq.client.CreateChatCompletion(s.ctx, req)
	if err != nil {
		return "", fmt.Errorf("groq title failed: %w", err)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func (s *Session) NewPrompt(message string) {
	// if len(s.History) < 1 {
	// 	s.nameSession(message)
	// }

	s.History = append(s.History, openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleUser,
		Content:    message,
	})

	// resp, err := s.groq.handlePrompt(s.ctx, s.History)
	// if err != nil {
	// 	s.History =  append(s.History, openai.ChatCompletionMessage{
	// 		Content: fmt.Sprintf("%v", err),
	// 	})
	// }

	go func() {
		time.Sleep(4*time.Second)
		resp := openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleAssistant,
			Content:    "Hello There!",
		}		
		s.History = append(s.History, resp)

		time.Sleep(2*time.Second)
		resp = openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleAssistant,
			Content:    "How are you?",
		}		
		s.History = append(s.History, resp)
	}()	
}

func InitSession() {
	CSession = Session{
		groq: newGroq(),
	}
}