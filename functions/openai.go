package functions

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func ChatGpt(params FunctionParams) (FunctionResponse, FunctionError) {
	requiredParams := []string{"prompt", "account"}
	for _, param := range requiredParams {
		if _, ok := params.Parameters[param]; !ok {
			return FunctionResponse{}, FunctionError{
				FunctionName: params.FunctionName,
				Message:      fmt.Sprintf("%s is required", param),
				Parameters:   params.Parameters,
				Trace:        fmt.Sprintf("%+v", params.Parameters),
			}
		}
	}

	// get the url, method, headers, and body from the params
	apiKey := params.Parameters["account"].(FunctionParams).Parameters["apiKey"].(string)
	prompt := params.Parameters["prompt"].(string)
	// Create a new HTTP request

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
			Trace:        fmt.Sprintf("%+v", err),
		}

	}

	// Return the response
	return FunctionResponse{
		FunctionName: params.FunctionName,
		Value: map[string]interface{}{
			"response": resp.Choices[0].Message.Content,
		},
		Parameters: params.Parameters,
	}, FunctionError{}
}
