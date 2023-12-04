package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// Example usage
//
//	emailParams := functions.FunctionParams{
//		FunctionName: "SendEmail",
//		Parameters: map[string]interface{}{
//			"apiKey":  "xxxxx",
//			"to":      "myestery@mailinator.com",
//			"subject": "Hello, Mailgun!",
//			"body":    "This is the body of the email.",
//			"domain":  "mg.xx.com",
//			"sender":  "jp@xx.com",
//		},
//	}
func SendEmail(params FunctionParams) (FunctionResponse, FunctionError) {
	// Validate required parameters
	requiredParams := []string{"apiKey", "to", "subject", "body", "domain", "sender"}
	for _, param := range requiredParams {
		if _, ok := params.Parameters[param]; !ok {
			return FunctionResponse{}, FunctionError{
				FunctionName: params.FunctionName,
				Message:      fmt.Sprintf("%s is required", param),
			}
		}
	}

	// Extract email parameters from the FunctionParams struct
	apiKey := params.Parameters["apiKey"].(string)
	to := params.Parameters["to"].(string)
	subject := params.Parameters["subject"].(string)
	body := params.Parameters["body"].(string)
	domain := params.Parameters["domain"].(string)
	sender := params.Parameters["sender"].(string)

	// Mailgun API endpoint and domain
	mg := mailgun.NewMailgun(domain, apiKey)
	message := mg.NewMessage(sender, subject, body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
			Trace:        fmt.Sprintf("%+v", err),
		}
	}

	// Return a response with status code and response headers
	return FunctionResponse{
		FunctionName: params.FunctionName,
		Value: map[string]interface{}{
			"id":       id,
			"response": resp,
		},
	}, FunctionError{}
}
