package gpt

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type GPTHandler struct {
	token string
}

func NewGPTHandler(token string) *GPTHandler {
	return &GPTHandler{
		token: token,
	}
}

func (h *GPTHandler) GetAnswerFromChat(prompt string) (string, error) {
	c := openai.NewClient(h.token)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	response, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion Error: %v", err)
		return "", err
	}

	prediction := response.Choices[0].Message.Content

	return prediction, nil
}

func (h *GPTHandler) CompleteText(prompt string, request openai.CompletionRequest) (string, error) {
	c := openai.NewClient(h.token)
	ctx := context.Background()
	var req openai.CompletionRequest

	if request.Model == "" {
		req = openai.CompletionRequest{
			Model:       openai.GPT3TextDavinci003,
			Prompt:      prompt,
			MaxTokens:   5,
			Temperature: 0.9,
			TopP:        1,
			N:           1,
		}
	} else {
		req = request
	}

	response, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion Error: %v", err)
		return "", err
	}

	prediction := response.Choices[0].Text

	fmt.Println("Choices :: ", response.Choices)

	return prediction, nil
}
