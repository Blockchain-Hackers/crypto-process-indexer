package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Error  interface{}            `json:"error"`
	Status bool                   `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func CCIP(params FunctionParams) (FunctionResponse, FunctionError) {

	requiredParams := []string{"account", "sourceChain", "destinationChain", "destinationAccount", "tokenAddress", "amount", "feeTokenAddress"}
	for _, param := range requiredParams {
		if _, ok := params.Parameters[param]; !ok {
			return FunctionResponse{}, FunctionError{
				FunctionName: params.FunctionName,
				Message:      fmt.Sprintf("%s is required", param),
			}
		}
	}

	// Extract email parameters from the FunctionParams struct
	// fmt.Println("account", params.Parameters["account"])
	// account is of Type FunctionParam, so we need to extract the apiKey from it
	privateKey := params.Parameters["account"].(FunctionParams).Parameters["privateKey"].(string)
	sourceChain := params.Parameters["sourceChain"].(string)
	destinationChain := params.Parameters["destinationChain"].(string)
	destinationAccount := params.Parameters["destinationAccount"].(string)
	tokenAddress := params.Parameters["tokenAddress"].(string)
	amount := params.Parameters["amount"].(string)
	feeTokenAddress := params.Parameters["feeTokenAddress"].(string)

	// send a http request to localhost 4444
	// {
	// 	"sourceChain": "ethereumSepolia",
	// 	"destinationChain": "polygonMumbai",
	// 	"destinationAccount": "0x9d087fC03ae39b088326b67fA3C788236645b717",
	// 	"tokenAddress": "0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05",
	// 	"amount": "10",
	// 	"privateKey": "6e1b485777de659f004d1133e422def4be77d0346716e65278a369c9eb9d544b",
	// 	"feeTokenAddress": null
	// }

	var url = "http://localhost:4444/ccip"
	var jsonStr = []byte(fmt.Sprintf(`{
		"sourceChain": "%s",
		"destinationChain": "%s",
		"destinationAccount": "%s",
		"tokenAddress": "%s",
		"amount": "%s",
		"privateKey": "%s",
		"feeTokenAddress": "%s"
	}`, sourceChain, destinationChain, destinationAccount, tokenAddress, amount, privateKey, feeTokenAddress))
	// fmt.Println("jsonStr", string(jsonStr))
	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)

	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
			Trace:        fmt.Sprintf("%+v", err),
		}
	}

	// check for status of not 200
	if resp.StatusCode != 200 {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			// always expect json message from ccip api asin resp.Body.message
			Message: fmt.Sprintf("CCIP API returned status code %d: %s", resp.StatusCode, resp.Body),
		}
	}

	var response Response
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
			Trace:        fmt.Sprintf("%+v", err),
		}
	}
	// Unmarshal JSON into the Response struct
	err = json.Unmarshal([]byte(respBody), &response)

	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
			Trace:        fmt.Sprintf("%+v", err),
		}
	}

	return FunctionResponse{
		FunctionName: params.FunctionName,
		// put resp.data in value
		Value:      response.Data,
		Parameters: params.Parameters,
	}, FunctionError{}
}
