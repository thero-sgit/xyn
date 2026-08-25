package ai

import (
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