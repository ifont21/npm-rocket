package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
	commitMessages := fmt.Sprintf("commit message:\n%s\n", commits)
	prompt := fmt.Sprintf(`I am working on a new version of an NPM library and need to create a CHANGELOG based on specific commit messages. 
	The changelog needs to be categorized into sections like Breaking Changes, Bug fixes, Features, Deprecated, and Refactor. 
	Please DO NOT include categories if there are no relevant commit messages.
	Here are the commit messages:
    %s
	the change log should follow the following format:
	
	- begin

	### Category

	| description                           | Commit Number |
	| ---------------------------------------- | ------------- |
	| [commit-message] | #124 |

	- end
	`, commitMessages)

	messages := []string{
		prompt,
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
			Temperature: 0.4,
			Messages:    chatMessages,
		})

		if err != nil {
			if strings.Contains(err.Error(), "401") {
				fmt.Printf("Error trying to generate the changelog, please check your openai token\n")
			} else {
				fmt.Printf("Completion Error: %v", err.Error())
			}
			break
		}
		finalResponse = response.Choices[0].Message.Content
	}
	/* fmt.Printf("Changelog response ************************* \n%s\n", finalResponse) */

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

	e := &openai.APIError{}
	if errors.As(err, &e) {
		switch e.StatusCode {
		case 401:
			return "", errors.New("invalid auth token or key")
		case 500:
			return "", errors.New("internal server error")
		case 429:
			fmt.Println("Too many request, trying again in 2 seconds ...")
			time.Sleep(2 * time.Second)
			response, err = c.CreateChatCompletion(ctx, req)
			if err != nil {
				return "", err
			}
		}
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
