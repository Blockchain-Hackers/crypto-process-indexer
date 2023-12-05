package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blockchain-hackers/indexer/functions"
	"github.com/blockchain-hackers/indexer/triggers"
	// "github.com/ethereum/go-ethereum/common"
)

func main() {
	resp, _ := callFunc("transfer", functions.FunctionParams{
		FunctionName: "transfer",
		Parameters: map[string]interface{}{
			"to":         "0x84188bc94B497131d0Ee2Cf7C154b22357c25208",
			"amount":     int64(500000000000000),
			"privateKey": "6e1b485777de659f004d1133e422def4be77d0346716e65278a369c9eb9d544b",
		}})

	fmt.Println(resp)
	triggers.Run()

	// Simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Hacker!")
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
