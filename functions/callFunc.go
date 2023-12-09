// callfunc is a function that takes a string functionName and returns the value of that function if it exists on this program
package functions

import (
	"time"

	"github.com/blockchain-hackers/indexer/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func ConvertDBParamsToFunctionParams(dbParams []database.Parameter, functionName string) FunctionParams {
	params := FunctionParams{
		Parameters: map[string]interface{}{},
	}
	for _, param := range dbParams {
		params.Parameters[param.Name] = param.Value
	}
	params.FunctionName = functionName
	return params
}

func ConvertFunctionParamsToDBParams(functionParams FunctionParams) []database.Parameter {
	params := []database.Parameter{}
	for key, value := range functionParams.Parameters {
		params = append(params, database.Parameter{
			Name:  key,
			Value: value,
		})
	}
	return params
}

func ConvertFunctionResponseToDBStep(functionResponse FunctionResponse) database.StepRun {
	return database.StepRun{
		Name: functionResponse.FunctionName + " on " + time.Now().Format(time.RFC3339),
		Parameters: ConvertFunctionParamsToDBParams(FunctionParams{
			Parameters: functionResponse.Parameters,
		}),
		Function: functionResponse.FunctionName,
		Logs:     functionResponse.Logs,
		Success:  true,
		Message:  functionResponse.Message,
		Value:    functionResponse.Value,
		ID:       primitive.NewObjectID(),
	}
}

func ConvertFunctionErrorToDBStep(functionError FunctionError) database.StepRun {
	return database.StepRun{
		Name: functionError.FunctionName + " on " + time.Now().Format(time.RFC3339),
		Parameters: ConvertFunctionParamsToDBParams(FunctionParams{
			Parameters: functionError.Parameters,
		}),
		Function: functionError.FunctionName,
		Logs:     functionError.Trace,
		Success:  false,
		Message:  functionError.Message,
		ID:       primitive.NewObjectID(),
	}
}
