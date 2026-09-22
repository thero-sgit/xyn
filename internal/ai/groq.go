package ai

import (
	"context"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/thero-sgit/xyn/internal/config"
)

type GroqChatCompletionMessage struct {
	Role 	string `json:"role"`
	Content string `json:"content"`
}

type GroqStreamRequest struct {
	Model 			string 						`json:"model"`
	Messages 		[]GroqChatCompletionMessage `json:"messages"`
	ReasoningFormat string 						`json:"reasoning_format,omitempty"`
	ReasoningEffort string                      `json:"reasoning_effort,omitempty"`
	Stream          bool   						`json:"stream"`
}

type GroqStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content 		 string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			Reasoning 		 string `json:"reasoning"`
		} `json:"delta"`
	} `json:"choices"`
}

type GroqChatCompletionStream struct {
	body   io.ReadCloser
	reader *bufio.Reader
	
}

func (s *GroqChatCompletionStream) Close() {
	s.body.Close()
}

func (s *GroqChatCompletionStream) Recv() (GroqStreamChunk, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return GroqStreamChunk{}, err
		}
	}

	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "data: ") {
		return GroqStreamChunk{}, nil
	}

	data := strings.TrimPrefix(line, "data: ")

	if data == "[DONE]" {
		return GroqStreamChunk{}, nil
	}

	var chunk GroqStreamChunk
	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return GroqStreamChunk{}, err
	}

	return chunk, nil
}

type groq struct {
	url    string
	apiKey string
}

func newGroq() *groq {
	return &groq{
		url:     "https://api.groq.com/openai/v1/chat/completions",
		apiKey:	 config.Config.GroqApiKey,
	}
}

func (g groq) CreateChatCompletionStream(ctx context.Context, reqBody GroqStreamRequest) (*GroqChatCompletionStream, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return &GroqChatCompletionStream{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", g.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &GroqChatCompletionStream{}, err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Initial request error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &GroqChatCompletionStream{}, fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(body))
	}

	return &GroqChatCompletionStream{
		body: resp.Body,
		reader: bufio.NewReader(resp.Body),
	}, nil
}

func (g groq) handlePrompt(ctx context.Context, messages []GroqChatCompletionMessage) (GroqChatCompletionStream, error) {
	req := GroqStreamRequest{
		Model: config.Config.Model,
		Messages:  messages,
		ReasoningEffort: "high",
		ReasoningFormat: "parsed",
		Stream: true,
	}

	stream, err := g.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return GroqChatCompletionStream{}, fmt.Errorf("groq completion failed: %w", err)
	}

	return *stream, nil
}
