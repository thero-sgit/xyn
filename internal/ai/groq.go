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

func (g groq) handlePrompt(ctx context.Context, messages []openai.ChatCompletionMessage) (openai.ChatCompletionStream, error) {
	req := openai.ChatCompletionRequest{
		Model: config.Config.Model,
		Messages:  messages,
	}

	stream, err := g.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionStream{}, fmt.Errorf("groq completion failed: %w", err)
	}

	return *stream, nil
}