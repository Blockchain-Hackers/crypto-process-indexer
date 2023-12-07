// callfunc is a function that takes a string functionName and returns the value of that function if it exists on this program
package functions

import (
	"github.com/blockchain-hackers/indexer/database"
)

// "github.com/blockchain-hackers/indexer/

func CallFunc(functionName string, param FunctionParams) (FunctionResponse, FunctionError) {
	var maps = map[string]func(_param FunctionParams) (FunctionResponse, FunctionError){
		"Transfer":          Transfer,
		"send-http-request": HTTPRequest,
		"send-email":        SendEmail,
	}

	if val, ok := maps[functionName]; ok {
		return val(param)
	}

	// no function was matched
	return FunctionResponse{}, FunctionError{}
}

func ConvertDBParamsToFunctionParams(dbParams []database.Parameter) FunctionParams {
	params := FunctionParams{
		Parameters: map[string]interface{}{},
	}
	for _, param := range dbParams {
		params.Parameters[param.Name] = param.Value
	}
	return params
}
