package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var Session session

type session struct {
	groq *groq
	name string	
}

func (s *session) nameSession(ctx context.Context, g groq, message string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: "llama-3.1-8b-instant",
		Messages:  []openai.ChatCompletionMessage{
			{Role: "system", Content: "Generate a 3-5 word title for this prompt. Return ONLY the title."},
			{Role: "user", Content: message},
		},
		MaxTokens: 15,
	}

	resp, err := g.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("groq title failed: %w", err)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func InitSession() {
	Session = session{
		groq: newGroq(),
	}
}