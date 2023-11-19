// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// // Payload struct represents the structure of the JSON payload
// type Payload struct {
// 	Level      string                 `json:"level"`
// 	Message    string                 `json:"message"`
// 	ResourceID string                 `json:"resourceId"`
// 	Timestamp  string                 `json:"timestamp"`
// 	TraceID    string                 `json:"traceId"`
// 	SpanID     string                 `json:"spanId"`
// 	Commit     string                 `json:"commit"`
// 	Metadata   map[string]interface{} `json:"metadata"`
// }

// // LogLevels is a slice of available log levels
// var LogLevels = []string{"info", "warn", "error", "debug"}

// // generateDummyData generates dummy data with random values
// func generateDummyData() Payload {
// 	rand.Seed(time.Now().UnixNano())

// 	return Payload{
// 		Level:      getRandomLogLevel(),
// 		Message:    "Failed to connect to DB",
// 		ResourceID: fmt.Sprintf("server-%d", rand.Intn(1000)),
// 		Timestamp:  time.Now().UTC().Format(time.RFC3339),
// 		TraceID:    getRandomID(),
// 		SpanID:     getRandomID(),
// 		Commit:     getRandomID(),
// 		Metadata: map[string]interface{}{
// 			"parentResourceId": fmt.Sprintf("server-%d", rand.Intn(1000)),
// 		},
// 	}
// }

// // getRandomID generates a random string ID
// func getRandomID() string {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	result := make([]byte, 10)
// 	for i := range result {
// 		result[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return string(result)
// }

// // getRandomLogLevel returns a random log level from LogLevels
// func getRandomLogLevel() string {
// 	return LogLevels[rand.Intn(len(LogLevels))]
// }
// func main() {
// 	othermain()
// }

// // func main() {

// // 	var wg sync.WaitGroup

// // 	// Number of requests to send
// // 	numRequests := 2000
// // 	wg.Add(numRequests)

// // 	var success int
// // 	// Retry up to 3 times
// // 	for i := 0; i < numRequests; i++ {
// // 		// var err error
// // 		var responseStatusCode int
// // 		for attempt := 1; attempt <= 3; attempt++ {
// // 			// Generate dummy data
// // 			payload := generateDummyData()

// // 			// Convert the payload to JSON
// // 			payloadBytes, err := json.Marshal(payload)
// // 			if err != nil {
// // 				fmt.Println("Error marshalling JSON:", err)
// // 				return
// // 			}

// // 			// Create a buffer with the JSON payload
// // 			buffer := bytes.NewBuffer(payloadBytes)

// // 			// Send the POST request with a timeout
// // 			url := "http://localhost:3000"
// // 			client := &http.Client{
// // 				Timeout: 5 * time.Second,
// // 			}

// // 			resp, err := client.Post(url, "application/json", buffer)
// // 			if err != nil {
// // 				fmt.Printf("Attempt %d - Error sending POST request: %s\n", attempt, err)
// // 				time.Sleep(2 * time.Second) // Wait before retrying
// // 				continue
// // 			}
// // 			defer resp.Body.Close()

// // 			// Record the response status code
// // 			responseStatusCode = resp.StatusCode

// // 			// If the request was successful (status code 2xx), break out of the retry loop
// // 			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
// // 				success++
// // 				break
// // 			} else {
// // 				fmt.Printf("Attempt %d - Non-successful response: %s\n", attempt, resp.Status)
// // 				time.Sleep(2 * time.Second) // Wait before retrying
// // 			}
// // 		}

// // 		// Print the final response status
// // 		fmt.Printf("Request %d - Final Response Status: %d\n", i, responseStatusCode)

// // 		// Mark the request as done
// // 		wg.Done()
// // 	}

// // 	wg.Wait()

// // 	fmt.Println("Successful requests: ", success, "out of", numRequests)
// // }

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Payload struct represents the structure of the JSON payload
type Payload struct {
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	ResourceID string                 `json:"resourceId"`
	Timestamp  string                 `json:"timestamp"`
	TraceID    string                 `json:"traceId"`
	SpanID     string                 `json:"spanId"`
	Commit     string                 `json:"commit"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// LogLevels is a slice of available log levels
var LogLevels = []string{"info", "warn", "error", "debug"}

// generateDummyData generates dummy data with random values
func generateDummyData() Payload {
	rand.Seed(time.Now().UnixNano())

	return Payload{
		Level:      getRandomLogLevel(),
		Message:    "Failed to connect to DB",
		ResourceID: fmt.Sprintf("server-%d", rand.Intn(1000)),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		TraceID:    getRandomID(),
		SpanID:     getRandomID(),
		Commit:     getRandomID(),
		Metadata: map[string]interface{}{
			"parentResourceId": fmt.Sprintf("server-%d", rand.Intn(1000)),
		},
	}
}

// getRandomID generates a random string ID
func getRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// getRandomLogLevel returns a random log level from LogLevels
func getRandomLogLevel() string {
	return LogLevels[rand.Intn(len(LogLevels))]
}

func main() {
	var wg sync.WaitGroup

	// Number of requests to send
	numRequests := 2000
	wg.Add(numRequests)

	var success int

	// Channel to communicate success or failure of each request
	resultChannel := make(chan int, numRequests)

	// Retry up to 3 times
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer wg.Done()

			var responseStatusCode int

			for attempt := 1; attempt <= 3; attempt++ {
				// Generate dummy data
				payload := generateDummyData()

				// Convert the payload to JSON
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					return
				}

				// Create a buffer with the JSON payload
				buffer := bytes.NewBuffer(payloadBytes)

				// Send the POST request with a timeout
				url := "http://localhost:3000"
				client := &http.Client{
					Timeout: 5 * time.Second,
				}

				resp, err := client.Post(url, "application/json", buffer)
				if err != nil {
					fmt.Printf("Attempt %d - Error sending POST request: %s\n", attempt, err)
					time.Sleep(2 * time.Second) // Wait before retrying
					continue
				}
				defer resp.Body.Close()

				// Record the response status code
				responseStatusCode = resp.StatusCode

				// If the request was successful (status code 2xx), break out of the retry loop
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					success++
					break
				} else {
					fmt.Printf("Attempt %d - Non-successful response: %s\n", attempt, resp.Status)
					time.Sleep(2 * time.Second) // Wait before retrying
				}
			}

			// Send the result to the result channel
			resultChannel <- responseStatusCode
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the result channel to signal that no more results will be sent
	close(resultChannel)

	// Process results from the channel
	for result := range resultChannel {
		fmt.Printf("Final Response Status: %d\n", result)
	}

	fmt.Println("Successful requests:", success, "out of", numRequests)
}
