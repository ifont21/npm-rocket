package adapters

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type GPTSuggestions struct {
	OpenAIToken string
}

func NewGPTSuggestion(openAIToken string) *GPTSuggestions {
	return &GPTSuggestions{
		OpenAIToken: openAIToken,
	}
}

func (g GPTSuggestions) GetSuggestedChangelogOutOfCommits(commits string) (string, error) {
	c := openai.NewClient(g.OpenAIToken)
	commit := fmt.Sprintf("commit message:\n%s\n", commits)
	mainPrompt := "can you generate a changelog based on these commits:"
	mainRequest := fmt.Sprintf("%s\n%s", mainPrompt, commit)
	messages := []string{
		mainRequest,
		`Can you consider now using this template for the changelog
		### Breaking changes
		-
		### Added
		-
		### Bug fixes
		-
		### Refactored
		-
		### Deprecated
		-`,
		"Now whenever you see the word `None` don't include the section",
		"Now add the commit messages in a markdown table with columns `commit message` and `description`. Add tables as many sections listed",
	}
	chatMessages := make([]openai.ChatCompletionMessage, 0)
	finalResponse := ""

	for _, message := range messages {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})
		response, err := c.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 0.2,
			Messages:    chatMessages,
		})
		if err != nil {
			fmt.Printf("Completion Error: %v", err)
			break
		}
		fmt.Println("Response: ", response.Choices[0].Message.Content)
		finalResponse = response.Choices[0].Message.Content
	}

	return finalResponse, nil
}

func (g GPTSuggestions) GetBumpTypeSuggestionOutOfCommits(commits string) (string, error) {
	c := openai.NewClient(g.OpenAIToken)
	request := `based on the commit message how you suggest to bump the library
		- major\n
		- minor\n
		- patch
	`
	prompt := fmt.Sprintf("%s\n%s", request, commits)

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

func (g GPTSuggestions) GetFilteredCommitsByScope(commits string, scope string) (string, error) {
	c := openai.NewClient(g.OpenAIToken)
	prompt := fmt.Sprintf("Select only the text related to `%s` based on the following commits, if there are not text related just say \"no commits related\"\n\n%s", scope, commits)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0.12,
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
