package main

import (
	// "context"
	"fmt"
	"log"
	"net/http"

	"github.com/blockchain-hackers/indexer/database"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	// "github.com/wailsapp/wails/lib/interfaces"
	// "go.mongodb.org/mongo-driver/bson"
	// "github.com/blockchain-hackers/indexer/functions"
	"github.com/blockchain-hackers/indexer/triggers"
	// "github.com/ethereum/go-ethereum/common"
)

func main() {

	DBClient := database.Connect()
	// dbs,_ := DBClient.ListDatabaseNames(context.Background(), bson.D{{}})
	fmt.Printf("Database connected successfully: %+v\n", (DBClient.NumberSessionsInProgress()))
	triggers.Run()
	// id, _ := primitive.ObjectIDFromHex("6576656fe4783994d0f6678d")
	// acc1, _ := database.GetAccount(id)

	// fmt.Printf("Account: %+v\n", acc1)

	// Simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Hacker!")
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
