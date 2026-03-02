package ai

import (
	"context"
)

type CompletionGenerator interface {
	GenerateCompletion(ctx context.Context, systemPrompt, prompt string) (string, error)
}

type AI struct {
	gen CompletionGenerator
}

func NewAi(completionGenerator CompletionGenerator) *AI {
	return &AI{
		gen: completionGenerator,
	}
}

func (ai *AI) GenWelcomePrompt(ctx context.Context, ownUserProfile, userProfile string) (string, error) {
	prompt := createWelcomePrompt(ownUserProfile, userProfile)
	output, err := ai.gen.GenerateCompletion(ctx, "", prompt)
	if err != nil {
		return "", err
	}

	return output, nil
}
