package ai

import (
    "context"
    "errors"
    "io"
    "sync"
    "time"

    "github.com/sashabaranov/go-openai"
)

type Chunk struct {
    Data string
    EOS  bool
}

var CSession Session

type Session struct {
    Ctx    context.Context
    cancl  context.CancelFunc
    groq   *groq
    Wg     sync.WaitGroup

    History []openai.ChatCompletionMessage
}

func (s *Session) SetNewContext() {
    if s.cancl != nil {
        s.cancl()
    }

    ctx, cancel := context.WithCancel(context.Background())
    s.Ctx = ctx
    s.cancl = cancel
}

func (s *Session) Close() {
    if s.cancl != nil {
        s.cancl()
    }
    s.Wg.Wait()
}

func (s *Session) NameSession(message string, c chan Chunk, errChan chan error) {
    req := openai.ChatCompletionRequest{
        Model: "openai/gpt-oss-20b",
        Messages: []openai.ChatCompletionMessage{
            {Role: "system", Content: "Generate a 3-5 word title for this prompt. Return ONLY the title. MUST BE 3 TO 5 WORDS!"},
            {Role: "user", Content: message},
        },
        Stream: true,
    }

    s.Wg.Add(1)
    go func() {
        defer s.Wg.Done()
        defer close(c)

        requestCtx, cancel := context.WithCancel(s.Ctx)
        defer cancel()

        const idleTimeout = 15 * time.Second
        timer := time.AfterFunc(idleTimeout, cancel)
        defer timer.Stop()

        stream, err := s.groq.client.CreateChatCompletionStream(requestCtx, req)
        if err != nil {
            s.handleError(err, errChan)
            return
        }
        defer stream.Close()

        for {
            select {
            case <-requestCtx.Done():
                return
            default:
            }

            response, err := stream.Recv()
            if err != nil {
                if errors.Is(err, io.EOF) {
                    c <- Chunk{Data: "", EOS: true}
                    return
                }
                s.handleError(err, errChan)
                return
            }

            if len(response.Choices) > 0 {
                c <- Chunk{Data: response.Choices[0].Delta.Content, EOS: false}
            }

            timer.Reset(idleTimeout)
        }
    }()
}

func (s *Session) NewPrompt(message string, c chan Chunk, errChan chan error) {
    h := append(s.History, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: message,
    })

    s.Wg.Add(1)
    go func() {
        defer s.Wg.Done()
        defer close(c)

        requestCtx, cancel := context.WithCancel(s.Ctx)
        defer cancel()

        const idleTimeout = 15 * time.Second
        timer := time.AfterFunc(idleTimeout, cancel)
        defer timer.Stop()

        var responseBuffer string

        stream, err := s.groq.handlePrompt(requestCtx, h)
        if err != nil {
            s.handleError(err, errChan)
            return
        }
        defer stream.Close()

        for {
            select {
            case <-requestCtx.Done():
                return
            default:
            }

            response, err := stream.Recv()
            if err != nil {
                if errors.Is(err, io.EOF) {
                    c <- Chunk{Data: "", EOS: true}
                    s.History = append(s.History, openai.ChatCompletionMessage{
                        Role:    openai.ChatMessageRoleUser,
                        Content: message,
                    })
                    s.History = append(s.History, openai.ChatCompletionMessage{
                        Role:    openai.ChatMessageRoleAssistant,
                        Content: responseBuffer,
                    })
                    return
                }
                s.handleError(err, errChan)
                return
            }

            timer.Reset(idleTimeout)

            if len(response.Choices) > 0 {
                content := response.Choices[0].Delta.Content
                c <- Chunk{Data: content, EOS: false}
                responseBuffer += content
            }
        }
    }()
}

func (s *Session) handleError(err error, errChan chan error) {
    if errChan != nil {
        select {
        case errChan <- err:
        default:
        }
    }
    if s.cancl != nil {
        s.cancl()
    }
}

func InitSession() {
    CSession = Session{
        groq: newGroq(),
    }
    CSession.SetNewContext()
}