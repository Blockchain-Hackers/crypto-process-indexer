// find a way to send request in go
// create a function that send post request with url and json


// http.Get()

// func callWebhook(){}


package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)


func makeHTTPRequest(url string,requestType string, requestBody []byte) bool {

	// Create an HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON payload
	req, err := http.NewRequest(requestType, url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Set the content type header for JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return false
	}

	// Print the response body
	fmt.Println("Response:", string(body))

	return true // Success
}


