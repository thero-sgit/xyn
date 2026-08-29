package ai

import (
	"context"
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

func (s *Session) NameSession(message string, c *chan string) {

	go func() {
		req := openai.ChatCompletionRequest{
			Model: "openai/gpt-oss-20b",
			Messages:  []openai.ChatCompletionMessage{
				{Role: "system", Content: "Generate a 3-5 word title for this prompt. Return ONLY the title."},
				{Role: "user", Content: message},
			},
		}

		resp, err := s.groq.client.CreateChatCompletion(s.ctx, req)
		if err != nil {
			panic(err.Error())
		}

		name := strings.TrimSpace(resp.Choices[0].Message.Content)

		*c <- name

	}()
}

func (s *Session) NewPrompt(message string, c *chan openai.ChatCompletionMessage) {
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

		*c <- resp

		time.Sleep(2*time.Second)
		resp = openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleAssistant,
			Content:    "How are you?",
		}		
		s.History = append(s.History, resp)

		*c <- resp
	}()	
}

func InitSession() {
	CSession = Session{
		ctx: context.Background(),
		groq: newGroq(),
	}
}