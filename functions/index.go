package functions

import (
	// "fmt"
	// "github.com/blockchain-hackers/indexer"
)

type ProcessFunctionParams struct {
	// The name of the function
	FunctionName string                 `json:"functionName"`
	Parameters   map[string]interface{} `json:"parameters"`
}

type ProcessFunctionError struct {
	// The name of the function
	FunctionName string `json:"functionName"`
	// The error message
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

type ProcessFunctionResponse struct {
	FunctionName string `json:"functionName"`
	Value        string `json:"value"`
}

