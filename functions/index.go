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
	Parameters map[string]interface{} `json:"parameters"`
}

func (e FunctionError) Exists() bool {
	return e.Message != ""
}

type FunctionResponse struct {
	FunctionName string                 `json:"functionName"`
	Value        map[string]interface{} `json:"value"`
	Parameters   map[string]interface{} `json:"parameters"`
	Logs         string               `json:"logs"`
	Message	  string                 `json:"message"`
}
