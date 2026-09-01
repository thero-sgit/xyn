package ai

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

type Chunk struct {
	Data string
	EOS  bool
	SOS  bool
}

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

func (s *Session) NameSession(message string, c chan Chunk, errChan chan error) {
	req := openai.ChatCompletionRequest{
		Model: "openai/gpt-oss-20b",
		Messages:  []openai.ChatCompletionMessage{
			{Role: "system", Content: "Generate a 3-5 word title for this prompt. Return ONLY the title."},
			{Role: "user", Content: message},
		},
		Stream: true,
	}

	go func() {	
		stream, err := s.groq.client.CreateChatCompletionStream(s.Ctx, req)
		if err != nil {
			s.handleError(err, errChan)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				c <- Chunk{Data: "", EOS: true}
				break
			}

			if err != nil {
				s.handleError(err, errChan)
			}

			if len(response.Choices) > 0 {
				c <- Chunk{Data: response.Choices[0].Delta.Content, EOS: false}
			}
		}
	}()
}

func (s *Session) NewPrompt(message string, c chan Chunk, errChan chan error) {
	s.History = append(s.History, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleUser,
			Content:    message,
	})

	go func() {
        stream, err := s.groq.handlePrompt(s.Ctx, s.History)
        if err != nil {
            s.handleError(err, errChan)
            return
        }
		
        for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				c <- Chunk{Data: "", EOS: true}
				break
			}

			if err != nil {
				s.handleError(err, errChan)
			}

			if len(response.Choices) > 0 {
				c <- Chunk{Data: response.Choices[0].Delta.Content, EOS: false}
			}
		}
    }()
}

func (s *Session) handleError(err error, errChan chan error) {
	s.cancl()
	errChan <- err
}

func InitSession() {
	CSession = Session{
		groq: newGroq(),
	}

	CSession.SetNewContext()
}