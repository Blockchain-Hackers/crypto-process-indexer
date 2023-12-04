package triggers
// package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainlinkPriceFeed struct {
}

var pairs map[string]string = map[string]string{
	"ETH/USD": "0x694AA1769357215DE4FAC081bf1f309aDC325306",
	"BTC/USD": "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43",
}

var ChainlinkContractABI = `[
	{
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "latestRoundData",
		"outputs": [
			{
				"internalType": "uint80",
				"name": "roundId",
				"type": "uint80"
			},
			{
				"internalType": "int256",
				"name": "answer",
				"type": "int256"
			},
			{
				"internalType": "uint256",
				"name": "startedAt",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "updatedAt",
				"type": "uint256"
			},
			{
				"internalType": "uint80",
				"name": "answeredInRound",
				"type": "uint80"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

// const infuraURL = "wss://sepolia.infura.io/ws/v3/927b0bef549145fba75661d347f23b8a"

// run triggers
func (trigger *ChainlinkPriceFeed) run() {
	for {
		for pair := range pairs {
			price, err := getLatestPrice(pair)
			if err != nil {
				log.Printf("Error getting latest price: %v", err)
				continue
			}
			var val, _ = json.Marshal(price)
			log.Printf("Latest price of %s: %s", pair, val)
			// wait for 15 seconds
		}
		time.Sleep(5 * time.Second)
	}
}

type ChainlinkLatestRoundData struct {
	Answer          *big.Int
	RoundId         *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func getLatestPrice(pair string) (*ChainlinkLatestRoundData, error) {

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(pairs[pair])
	parsedABI, err := abi.JSON(strings.NewReader(ChainlinkContractABI))
	// fmt.Printf("Parsed ABI: %v", parsedABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	query := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["latestRoundData"].ID,
	}

	result, err := client.CallContract(context.Background(), query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}
	// log result
	// log.Printf("Result: %v", result)

	// var answer big.Int
	var data, _ = parsedABI.Unpack("latestRoundData", result)
	log.Printf("Data: %v", data)
	var answer2 ChainlinkLatestRoundData
	err = parsedABI.UnpackIntoInterface(&answer2, "latestRoundData", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack contract result: %v", err)
	}

	return &answer2, nil
}

// func main() {
// 	trigger := &ChainlinkPriceFeed{}
// 	trigger.run()
// }
