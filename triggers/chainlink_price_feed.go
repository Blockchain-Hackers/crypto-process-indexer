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

	"github.com/blockchain-hackers/indexer/database"
	"github.com/blockchain-hackers/indexer/runner"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
)

type ChainlinkPriceFeed struct {
}

var pairs map[string]string = map[string]string{
	"ETH/USDT": "0x694AA1769357215DE4FAC081bf1f309aDC325306",
	"BTC/USDT": "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43",
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
	fmt.Println("Running Chainlink Price Feed trigger...")
	// fmt.Println(DBClient)
	for {
		for pair := range pairs {
			price, err := getLatestPrice(pair)
			if err != nil {
				log.Printf("Error getting latest price: %v", err)
				continue
			}
			// price.
			var val, _ = json.Marshal(price)
			log.Printf("Latest price of %s: %s", pair, val)
			// wait for 15 seconds
			go func(pair string, price *ChainlinkLatestRoundData) {
				var decimals int64 = 100000000
				var Price float64 = float64(price.Answer.Int64()) / float64(decimals)
				fmt.Println("Pair: ", pair)
				fmt.Println("Price: ", Price)
				trigger.processEvent(Event{
					EventName: "ListenForPriceChanges",
					Data: map[string]interface{}{
						"Pair":  pair,
						"Price": Price,
					},
				})
			}(pair, price)
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

func (trigger *ChainlinkPriceFeed) processEvent(event Event) {
	// look for all flows with this event name as triger using mongoDb and the event name
	// if found, run the flow
	// if not found, do nothing

	// we have to do for greater than, grater than or equal and also less than, less than or equal and then equal

	// so we first get those that wanted greater than or equal
	filter1 := bson.M{
		"trigger.name": event.EventName,
		"trigger.parameters": bson.M{
			"$all": []bson.M{
				{"$elemMatch": bson.M{"name": "pair", "value": event.Data["Pair"]}},
				{"$elemMatch": bson.M{"name": "condition", "value": ">="}},
				{"$elemMatch": bson.M{"name": "value", "value": bson.M{"$lte": event.Data["Price"]}}},
			},
		},
	}
	// fmt.Printf("Filter: %+v\n", filter1)
	// then we get those that wanted greater than
	filter2 := bson.M{
		"trigger.name": event.EventName,
		"trigger.parameters": bson.M{
			"$all": []bson.M{
				{"$elemMatch": bson.M{"name": "pair", "value": event.Data["Pair"]}},
				{"$elemMatch": bson.M{"name": "condition", "value": ">"}},
				{"$elemMatch": bson.M{"name": "value", "value": bson.M{"$lt": event.Data["Price"]}}},
			},
		},
	}
	// then we get those that wanted less than or equal
	filter3 := bson.M{
		"trigger.name": event.EventName,
		"trigger.parameters": bson.M{
			"$all": []bson.M{
				{"$elemMatch": bson.M{"name": "pair", "value": event.Data["Pair"]}},
				{"$elemMatch": bson.M{"name": "condition", "value": "<="}},
				{"$elemMatch": bson.M{"name": "value", "value": bson.M{"$gte": event.Data["Price"]}}},
			},
		},
	}
	// then we get those that wanted less than
	filter4 := bson.M{
		"trigger.name": event.EventName,
		"trigger.parameters": bson.M{
			"$all": []bson.M{
				{"$elemMatch": bson.M{"name": "pair", "value": event.Data["Pair"]}},
				{"$elemMatch": bson.M{"name": "condition", "value": "<"}},
				{"$elemMatch": bson.M{"name": "value", "value": bson.M{"$gt": event.Data["Price"]}}},
			},
		},
	}
	// then we get those that wanted equal
	filter5 := bson.M{
		"trigger.name": event.EventName,
		"trigger.parameters": bson.M{
			"$all": []bson.M{
				{"$elemMatch": bson.M{"name": "pair", "value": event.Data["Pair"]}},
				{"$elemMatch": bson.M{"name": "condition", "value": "="}},
				{"$elemMatch": bson.M{"name": "value", "value": bson.M{"$eq": event.Data["Price"]}}},
			},
		},
	}
	for _, filter := range []bson.M{filter1, filter2, filter3, filter4, filter5} {
		var flows = database.FindFlows(filter)
		fmt.Printf("Found %d flows\n", len(flows))

		for _, flow := range flows {
			fmt.Printf("Running flow: %+v\n", flow.Name)
			// get the steps and run them in series
			go runner.Run(flow)
		}
	}

}

func (trigger *ChainlinkPriceFeed) EventName() string {
	return "ChainlinkPriceFeed"
}
