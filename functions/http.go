package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HTTPRequest(params FunctionParams) (FunctionResponse, FunctionError) {
	// Create an HTTP client
	client := &http.Client{}

	requiredParams := []string{"url", "method", "headers", "body"}
	for _, param := range requiredParams {
		if _, ok := params.Parameters[param]; !ok {
			return FunctionResponse{}, FunctionError{
				FunctionName: params.FunctionName,
				Message:      fmt.Sprintf("%s is required", param),
			}
		}
	}

	// get the url, method, headers, and body from the params
	url := params.Parameters["url"].(string)
	method := params.Parameters["method"].(string)
	rawheaders := params.Parameters["headers"].(primitive.D)
	rawBody := params.Parameters["body"].(primitive.D)
	// var body map[string]interface{}
	body, _ := bson.Marshal(rawBody)
	headers, _ := bson.Marshal(rawheaders)
	bodyMap := map[string]interface{}{}
	headersMap := map[string]interface{}{}
	bson.Unmarshal(body, &bodyMap)
	bson.Unmarshal(headers, &headersMap)
	jsonBody, _ := json.Marshal(bodyMap)
	// Create a new HTTP request 
	req, err := http.NewRequest(method, url,  bytes.NewBuffer(jsonBody))
	// bytes.NewBuffer([]byte(body)))
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	// Add headers to the request
	for key, value := range headersMap {
		req.Header.Add(key, value.(string))
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
			Trace:        fmt.Sprintf("%+v", err),
		}
	}

	// Read the response body
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBody := buf.String()
	var respBodyJSON interface{}
	// parse the response body as JSON if the content-type is application/json
	if resp.Header.Get("Content-Type") == "application/json" {
		err = json.Unmarshal([]byte(respBody), &respBodyJSON)
		if err != nil {
			respBodyJSON = respBody
		}
	} else {
		respBodyJSON = respBody
	}

	// Return the response
	return FunctionResponse{
		FunctionName: params.FunctionName,
		Value: map[string]interface{}{
			"statusCode":      resp.StatusCode,
			"responseHeaders": resp.Header,
			"body":            respBodyJSON,
		},
	}, FunctionError{}
}
