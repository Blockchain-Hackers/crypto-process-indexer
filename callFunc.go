// callfunc is a function that takes a string functionName and returns the value of that function if it exists on this program
package main

import (
	"github.com/blockchain-hackers/indexer/functions"
)

func callFunc(functionName string, param functions.FunctionParams) (functions.FunctionResponse, functions.FunctionError) {
	var maps = map[string]func(_param functions.FunctionParams) (functions.FunctionResponse, functions.FunctionError){
		"Transfer": functions.Transfer,
	}

	if val, ok := maps[functionName]; ok {
		return val(param)
	}

	// no function was matched
	return functions.FunctionResponse{      
	}, functions.FunctionError{}
}
