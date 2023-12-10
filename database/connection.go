package database

// package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var FlowsCollection *mongo.Collection
var TriggersCollection *mongo.Collection

// return a MongoDB client
func Connect() *mongo.Client {
	_err := godotenv.Load(".env")
	if _err != nil {
		panic("Error loading .env file")
	}
	uri := os.Getenv("MONGO_URL")
	fmt.Println("Connecting to MongoDB...")
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your Database. You successfully connected to MongoDB!")
	Client = client
	FlowsCollection = Client.Database("cp").Collection("flows")
	TriggersCollection = Client.Database("cp").Collection("triggers")
	return client
}

func FindFlowsByTrigger(triggerName string) []Workflow {
	var cursor, err = FlowsCollection.Find(context.Background(), map[string]interface{}{
		"trigger.name": triggerName,
	})
	var flows []Workflow
	if err != nil {
		fmt.Println("Error: ", err)
		return flows
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var data Workflow
		err = cursor.Decode(&data)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			flows = append(flows, data)
		}
	}
	return flows
}

func FindFlows(filter bson.M) []Workflow {
	var cursor, err = FlowsCollection.Find(context.Background(), filter)
	var flows []Workflow
	if err != nil {
		fmt.Println("Error: ", err)
		return flows
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var data Workflow
		err = cursor.Decode(&data)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			flows = append(flows, data)
		}
	}
	return flows
}

func WriteRunToFlow(flowID primitive.ObjectID, run FlowRun) {
	// write the run to the flow runs collection
	_, err := FlowsCollection.UpdateOne(context.Background(), bson.M{"_id": flowID}, bson.M{
		"$push": bson.M{
			"runs": run,
		},
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func GetAccount(ID primitive.ObjectID) (Account, error) {
	var account Account
	err := Client.Database("cp").Collection("accounts").FindOne(context.Background(), bson.M{
		"_id": ID,
	}).Decode(&account)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}
