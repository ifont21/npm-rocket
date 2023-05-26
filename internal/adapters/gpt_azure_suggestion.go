package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type GPTAzureSuggestions struct {
	apiKey  string
	baseUrl string
}

func NewGPTAzureSuggestions(key string, baseUrl string) *GPTAzureSuggestions {
	return &GPTAzureSuggestions{
		apiKey:  key,
		baseUrl: baseUrl,
	}
}

func (g GPTAzureSuggestions) GetFilteredCommitsByScope(commits string, scope string) (string, error) {
	prompt := fmt.Sprintf("Select only the text related to `%s` based on the following commits, if there is no text related just say \"no commits related\"\n\n%s", scope, commits)
	config := openai.DefaultAzureConfig(g.apiKey, g.baseUrl, "gpt-4")
	client := openai.NewClientWithConfig(config)
	response, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			fmt.Printf("Error trying to filter the commits, please check your openai token\n")
		} else {
			fmt.Printf("Completion Error: %v", err)
		}
		return "", err
	}

	prediction := response.Choices[0].Message.Content

	return prediction, nil
}

func (g GPTAzureSuggestions) GetBumpTypeSuggestionOutOfCommits(commits string) (string, error) {
	request := `based on the commit message how you suggest to bump the library
		- major\n
		- minor\n
		- patch
	`
	prompt := fmt.Sprintf("%s\n%s", request, commits)

	config := openai.DefaultAzureConfig(g.apiKey, g.baseUrl, "gpt-4")
	client := openai.NewClientWithConfig(config)
	response, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			fmt.Printf("Error trying to filter the commits, please check your openai token\n")
		} else {
			fmt.Printf("Completion Error: %v", err)
		}
		return "", err
	}

	prediction := response.Choices[0].Message.Content

	return prediction, nil
}

func (g GPTAzureSuggestions) GetSuggestedChangelogOutOfCommits(commits string) (string, error) {
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
		"Now add the commit messages in a markdown table with columns `commit message` and `description`. Add tables as many sections listed and whenever you see the word `None` don't include the section",
	}
	chatMessages := make([]openai.ChatCompletionMessage, 0)
	finalResponse := ""

	config := openai.DefaultAzureConfig(g.apiKey, g.baseUrl, "gpt-4")
	client := openai.NewClientWithConfig(config)

	for _, message := range messages {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})
		response, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Messages: chatMessages,
		})
		if err != nil {
			if strings.Contains(err.Error(), "401") {
				fmt.Printf("Error trying to filter the commits, please check your openai token\n")
			} else {
				fmt.Printf("Completion Error: %v", err)
			}
			return "", err
		}
		finalResponse = response.Choices[0].Message.Content
		fmt.Printf("Changelog response ************************* \n%s\n", finalResponse)
	}
	return finalResponse, nil
}
