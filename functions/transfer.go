package functions

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Transfer(params FunctionParams) (FunctionResponse, FunctionError) {
	const infuraURL = "https://sepolia.infura.io/v3/927b0bef549145fba75661d347f23b8a"
	requiredParams := []string{"amount", "account", "recipient"}
	for _, param := range requiredParams {
		if _, ok := params.Parameters[param]; !ok {
			return FunctionResponse{}, FunctionError{
				FunctionName: params.FunctionName,
				Message:      fmt.Sprintf("%s is required", param),
			}
		}
	}

	// get the amount, to, and private key from the params
	amount := params.Parameters["amount"].(string)
	to := params.Parameters["recipient"].(string)
	// privateKey := params.Parameters["privateKey"].(string)
	privateKey := params.Parameters["account"].(FunctionParams).Parameters["privateKey"].(string)

	// create a new client
	client, err := ethclient.DialContext(context.Background(), infuraURL)
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	_privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	publicKey := _privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      "error casting public key to ECDSA",
		}
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// generate uuid as nonce

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}
	var intAmount, _ = strconv.Atoi(amount)
	tx := types.NewTransaction(nonce, common.HexToAddress(to), big.NewInt(
		int64(intAmount),
		// amount,
	), 2100000, gasPrice, nil)

	// Sign the transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), _privateKey)
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return FunctionResponse{}, FunctionError{
			FunctionName: params.FunctionName,
			Message:      err.Error(),
		}
	}

	return FunctionResponse{
		FunctionName: params.FunctionName,
		Value: map[string]interface{}{
			"txHash":   signedTx.Hash().Hex(),
			"gasLimit": signedTx.Gas(),
		},
		Parameters: params.Parameters,
	}, FunctionError{}
}
