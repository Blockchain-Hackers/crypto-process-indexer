// callfunc is a function that takes a string functionName and returns the value of that function if it exists on this program
package functions

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/blockchain-hackers/indexer/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// "github.com/blockchain-hackers/indexer/

func CallFunc(functionName string, param FunctionParams) (FunctionResponse, FunctionError) {
	var maps = map[string]func(_param FunctionParams) (FunctionResponse, FunctionError){
		"send-crypto":          Transfer,
		"send-http-request": HTTPRequest,
		"send-email":        SendEmail,
		"chat-gpt":          ChatGpt,
	}

	if val, ok := maps[functionName]; ok {
		return val(param)
	}

	// no function was matched
	return FunctionResponse{}, FunctionError{}
}

func ConvertDBParamsToFunctionParams(
	dbParams []database.Parameter,
	functionName string,
	triggerValue map[string]interface{},
	steps []database.StepRun,
	) FunctionParams {
	var ValuesMap = constructMap(triggerValue, steps)
	// now construct a s
	params := FunctionParams{
		Parameters: map[string]interface{}{},
	}
	for _, param := range dbParams {
		switch param.Type {
		case "string":
			params.Parameters[param.Name] = replacePlaceholders(param.Value.(string), ValuesMap)
		case "account":
			var accountID = string(param.Value.(string))
			var accountIDMongo, _ = primitive.ObjectIDFromHex(accountID)
			var resolvedAccount, _ = database.GetAccount(accountIDMongo)
			params.Parameters[param.Name] = ConvertDBParamsToFunctionParams(resolvedAccount.Parameters, functionName, triggerValue, []database.StepRun{})
			// fmt.Println("resolvedAccount: ", resolvedAccount)
			// fmt.Println("Param", param)
		default:
			params.Parameters[param.Name] = param.Value
		}

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

// StepsItem represents an item in the steps array.
type StepsItem struct {
	Value map[string]interface{}
}

// Steps represents the structure of the steps struct.
type Steps struct {
	Items []StepsItem
}

func constructMap(triggerValue map[string]interface{}, steps []database.StepRun) map[string]interface{} {
	result := make(map[string]interface{})
	// Add values from the triggerValue map
	for key, value := range triggerValue {
		result["flow.trigger.outputs."+key] = value
	}
	// Add values from the steps struct
	for i, step := range steps {
		for key, value := range step.Value {
			result[fmt.Sprintf("flow.steps[%d].value.%s", i, key)] = value
		}
	}
	return result
}

func replacePlaceholders(input string, values map[string]interface{}) string {
	re := regexp.MustCompile(`{{\s*([\w.]+)\s*}}`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the path from the placeholder
		path := strings.TrimSpace(strings.Trim(match, "{}"))

		// Get the value from the map using the path
		value, ok := values[path]
		if ok {
			// Replace the placeholder with the value
			return fmt.Sprintf("%v", value)
		}

		// If the path is not found in the map, keep the original placeholder
		return match
	})

	return result
}
