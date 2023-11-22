package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Replace 'YOUR_GOERLI_INFURA_API_KEY' with your Infura API key or use your own Ethereum node
const infuraURL = "wss://sepolia.infura.io/ws/v3/927b0bef549145fba75661d347f23b8a"

// Replace 'YOUR_CONTRACT_ADDRESSES' with an array of contract addresses you want to listen to
var contractAddresses = []string{"0xA17ddf0a5309d50D7a69CA096A5473240A715DfA"}

// Replace 'YOUR_CONTRACT_ABI' with the ABI (Application Binary Interface) of your contract
var contractAbi = `[{"inputs":[{"internalType":"string","name":"_greeting","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"greeting","type":"string"}],"name":"GreetingSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"action","type":"string"}],"name":"Interaction","type":"event"},{"inputs":[],"name":"performInteraction","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"_newGreeting","type":"string"}],"name":"setGreeting","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getGreeting","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

// create a map from contract address to contract ABI
var contractAbis = map[string]*abi.ABI{}

func main() {
	// Initialize Ethereum client
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// map contract addresses to lower case
	for i := range contractAddresses {
		contractAddresses[i] = strings.ToLower(contractAddresses[i])
	}

	// Create contract instances
	contracts := make([]*abi.ABI, len(contractAddresses))
	for i, _ := range contractAddresses {
		// contractAddress := common.HexToAddress(address)
		contractAbi, err := abi.JSON(strings.NewReader(contractAbi))
		if err != nil {
			log.Fatal(err)
		}
		contracts[i] = &contractAbi
		contractAbis[contractAddresses[i]] = &contractAbi
	}

	latestBlockNumber := uint64(0)

	// Subscribe to new block headers
	go func() {
		for {
			latestBlock, err := client.BlockByNumber(context.Background(), nil)
			if err != nil {
				log.Printf("Error getting latest block number: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if latestBlock.NumberU64() > latestBlockNumber {
				fmt.Printf("Latest block number: %d\n", latestBlock.NumberU64())
				getAllEventsInBlock(client, contracts, latestBlock.NumberU64())
				latestBlockNumber = latestBlock.NumberU64()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// Simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Hacker!")
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Get all events in a given block for specified contract addresses
func getAllEventsInBlock(client *ethclient.Client, contracts []*abi.ABI, blockNumber uint64) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Printf("Error getting block: %v", err)
		return
	}

	for _, tx := range block.Transactions() {
		fmt.Printf("Transaction: to %+v\n", tx.To())
		if tx.To() != nil && contains(contractAddresses, strings.ToLower(tx.To().Hex())) {
			fmt.Printf("New transaction on contract %s\n", tx.To().Hex())
			// fmt.Printf("Transaction: %+v\n", tx)

			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				log.Printf("Error getting transaction receipt: %v", err)
				continue
			}

			logs := receipt.Logs
			fmt.Printf("Logs: %+v\n", logs)

			for _, log := range logs {
				contract := findContract(contracts, strings.ToLower(log.Address.Hex()))
				if contract != nil {
					// Decode logs using the event ABI
					event, err := contract.EventByID(log.Topics[0])
					if err != nil {
						fmt.Printf("Error getting event by ID: %v", err)
						continue
					}

					// Unpack the log data
					// eventData := []interface{ new(interface{}) }
					// err, eventData =
					var eventData interface{}
					eventData, err = event.Inputs.Unpack(log.Data)
					if err != nil {
						fmt.Printf("Error unpacking log data: %v", err)
						continue
					}

					fmt.Printf("Decoded Logs for event %s: %+v\n", event.Name, eventData)
				} else {
					fmt.Printf("Contract ABI not found for address %s. Unable to decode event.", log.Address.Hex())
				}
			}
		}
	}
}

// Helper function to check if a string is in a slice of strings
func contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}

// Helper function to find a contract by address using the map
func findContract(contracts []*abi.ABI, address string) *abi.ABI {
	if contract, ok := contractAbis[address]; ok {
		return contract
	}
	return nil
}
