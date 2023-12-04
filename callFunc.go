// callfunc is a function that takes a string functionName and returns the value of that function if it exists on this program
package main

import (
	"github.com/blockchain-hackers/indexer/functions"
)

func callFunc(functionName string, param functions.ProcessFunctionParams) (functions.ProcessFunctionResponse, functions.ProcessFunctionError) {
	var maps = map[string]func(_param functions.ProcessFunctionParams) (functions.ProcessFunctionResponse, functions.ProcessFunctionError) {
		"Transfer": functions.Transfer,
	}

	if val, ok := maps[functionName]; ok {
		return val(param)
	}

	return functions.ProcessFunctionResponse{
		FunctionName: functionName,
		Value:        "",
	}, functions.ProcessFunctionError{}
}
