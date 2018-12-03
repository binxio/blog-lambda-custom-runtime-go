package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	awsLambdaRuntimeApi := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	if awsLambdaRuntimeApi == "" {
			panic("Missing: 'AWS_LAMBDA_RUNTIME_API'")
	}
	for {
		// get the next event
		requestUrl := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/next", awsLambdaRuntimeApi)
		resp, err := http.Get(requestUrl)
		if err != nil {
			log.Fatal(fmt.Errorf("Expected status code 200, got %d", resp.StatusCode))
		}

		requestId := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
		// print the next event
		eventData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(fmt.Errorf("Error: %s"), err)
		}
		fmt.Println("Received event:", string(eventData))

		// Assume API Gateway and respond with Hello World
		responseUrl := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/%s/response", awsLambdaRuntimeApi, requestId)
		responsePayload := []byte(`{"statusCode": 200, "body": "Hello World!"}`)
		req, err := http.NewRequest("POST", responseUrl, bytes.NewBuffer(responsePayload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		client.Timeout = time.Second * 1
		postResp, err := client.Do(req)
		if err != nil {
			log.Fatal(fmt.Errorf("Error %s", err))
		}
		body, _ := ioutil.ReadAll(postResp.Body)
		fmt.Println("Received response:", string(body))
	}
}