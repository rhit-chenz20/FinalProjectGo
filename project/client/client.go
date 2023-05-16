package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	stressTestBasicGet()
	stressTestBasicPost()
}

const AMOUNT_OF_REQ = 1000
const BASIC_SERVER = "http://127.0.0.1:8080"

// stress test - high volume of GET requests
func stressTestBasicGet() time.Duration {
	start := time.Now()
	for i := 0; i < AMOUNT_OF_REQ; i++ {
		resp, err := http.Get(BASIC_SERVER + "/students/")
		if err != nil {
			fmt.Println(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil || body == nil {
			fmt.Println(err)
		}
	}
	duration := time.Since(start)
	return duration
}

func stressTestBasicPost() time.Duration {
	start := time.Now()
	for i := 0; i < AMOUNT_OF_REQ; i++ {
		postBody, _ := json.Marshal(map[string]string{
			"ID":         "1",
			"coursename": "CSSE304",
			"grade":      fmt.Sprintf("%d", i),
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(BASIC_SERVER+"/students/", "application/json", responseBody)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil || body == nil {
			fmt.Println(err)
		}
	}
	duration := time.Since(start)
	return duration
}
