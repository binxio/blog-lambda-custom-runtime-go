package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	awsLambdaRuntimeApi := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	if awsLambdaRuntimeApi == "" {
			panic("Missing: 'AWS_LAMBDA_RUNTIME_API'")
	}
	for {
		// get the next event
		requestUrl := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/next", awsLambdaRuntimeApi)
		nextEventResponse, err := http.Get(requestUrl)
		HandleError(err)

		requestId := nextEventResponse.Header.Get("Lambda-Runtime-Aws-Request-Id")
		// print the next event
		eventData, err := ioutil.ReadAll(nextEventResponse.Body)
		HandleError(err)
		fmt.Println("Received event:", string(eventData))

		// Assume API Gateway and respond with Hello World
		responseUrl := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/%s/response", awsLambdaRuntimeApi, requestId)
		responsePayload := []byte(`{"statusCode": 200, "body": "Hello World!"}`)
		req, err := http.NewRequest("POST", responseUrl, bytes.NewBuffer(responsePayload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		client.Timeout = time.Second * 1
		resp, err := client.Do(req)
		HandleError(err)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Received response:", string(body))
	}
}