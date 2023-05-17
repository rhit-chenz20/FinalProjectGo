package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const NUM_WORKERS = 3
const AMOUNT_OF_REQ = 1000
const BASIC_SERVER_1 = "http://127.0.0.1:8080"
const BASIC_SERVER_2 = "http://127.0.0.1:8081"
const BASIC_SERVER_3 = "http://127.0.0.1:8002"

const CONCURRENT_SERVER_1 = "http://127.0.0.1:8083"
const CONCURRENT_SERVER_2 = "http://127.0.0.1:8084"
const CONCURRENT_SERVER_3 = "http://127.0.0.1:8085"

//three goroutines for each server that will spam requests
//each goroutine will have an output channel that the goroutine will send
//a message through once it is done
//once the main method received three messages through an output channel,
//end the timer.
func main() {
	//Stress test basic GET
	basic_time := time.Now()
	basic_c := make(chan int, NUM_WORKERS)
	go stressTestBasicGet(BASIC_SERVER_1, basic_c)
	go stressTestBasicGet(BASIC_SERVER_2, basic_c)
	go stressTestBasicGet(BASIC_SERVER_3, basic_c)

	for i := 0; i < NUM_WORKERS; i++ {
		val := <- basic_c
		fmt.Println("done with %d", val)
	}
	basic_duration := time.Since(basic_time)
	fmt.Println("Stress test GET requests basic server time: %d", basic_duration)

	//Stress test concurrent GET
	// concurrent_time := time.Now()
	// concurrent_c := make(chan int, NUM_WORKERS)
	// go stressTestBasicGet(CONCURRENT_SERVER_1, concurrent_c)
	// go stressTestBasicGet(CONCURRENT_SERVER_2, concurrent_c)
	// go stressTestBasicGet(CONCURRENT_SERVER_3, concurrent_c)

	// for i := 0; i < NUM_WORKERS; i++ {
	// 	val_c := <- concurrent_c
	// 	fmt.Println("done with %d", val_c)
	// }
	// concurrent_duration := time.Since(concurrent_time)
	// fmt.Println("Stress test GET requests concurrent server time: %d", concurrent_duration)
}



//--------------------------------------------------
// Basic Server (w locking) Tests
func stressTestBasicGet(port string, basic_c chan int) {
	for i := 0; i < AMOUNT_OF_REQ; i++ {
		resp, err := http.Get(port + "/students/")
		if err != nil {
			fmt.Println(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil || body == nil {
			fmt.Println(err)
		}
	}
	basic_c <- 1
}

func stressTestBasicPost(port string, basic_c chan int) {
	for i := 0; i < AMOUNT_OF_REQ; i++ {
		postBody, _ := json.Marshal(map[string]string{
			"ID":         "1",
			"coursename": "CSSE304",
			"grade":      fmt.Sprintf("%d", i),
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(port+"/students/", "application/json", responseBody)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil || body == nil {
			fmt.Println(err)
		}
	}
	basic_c <- 1
}
//--------------------------------------------------

//--------------------------------------------------
//Concurrent Server tests
func stressTestConcurrentGet(port string, concurrent_c chan int) {
	for i := 0; i < AMOUNT_OF_REQ; i++ {
		resp, err := http.Get(port + "/students/")
		if err != nil {
			fmt.Println(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil || body == nil {
			fmt.Println(err)
		}
	}
	concurrent_c <- 1 
}
//--------------------------------------------------