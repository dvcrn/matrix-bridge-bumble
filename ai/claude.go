package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
)

type Claude struct {
	anthropicClient *anthropic.Client
}

func NewClaude(apiKey string) *Claude {
	client := anthropic.NewClient(apiKey)

	return &Claude{
		anthropicClient: client,
	}
}

func (c *Claude) GenerateCompletion(ctx context.Context, systemPrompt, prompt string) (string, error) {
	resp, err := c.anthropicClient.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Dot5Sonnet20240620,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(prompt),
		},
		System:    systemPrompt,
		MaxTokens: 1000,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}

		return "", err
	}

	return resp.Content[0].GetText(), nil
}
