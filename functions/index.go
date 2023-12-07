package functions

// "fmt"
// "github.com/blockchain-hackers/indexer"

type FunctionParams struct {
	// The name of the function
	FunctionName string                 `json:"functionName"`
	Parameters   map[string]interface{} `json:"parameters"`
}

type FunctionError struct {
	// The name of the function
	FunctionName string `json:"functionName"`
	// The error message
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

func (e FunctionError) Exists() bool {
	return e.Message != ""
}

type FunctionResponse struct {
	FunctionName string                 `json:"functionName"`
	Value        map[string]interface{} `json:"value"`
}
