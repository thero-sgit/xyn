package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/thero-sgit/xyn/internal/config"
)

type groq struct {
	client *openai.Client
}

func newGroq() *groq {
	return &groq{
		client: openai.NewClientWithConfig(
			config.Config.GroqClientConfig,
		),
	}
}

func (g groq) handlePrompt(ctx context.Context, messages []openai.ChatCompletionMessage) (openai.ChatCompletionMessage, error) {
	req := openai.ChatCompletionRequest{
		Model: "llama-3.1-8b-instant",
		Messages:  messages,
	}

	resp, err := g.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, fmt.Errorf("groq title failed: %w", err)
	}

	return resp.Choices[0].Message, nil
}