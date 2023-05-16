package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
)

//stress test emulation of a system with high traffic volume
func main(){
	start := time.Now()
	for i := 0; i < 1000; i++ {
		resp, err := http.Get("http://127.0.0.1:8080/students/")
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil || body == nil {
			fmt.Println(err)
		}

		//print
		// sb := string(body)
		// fmt.Println(sb)
	}
	duration := time.Since(start)
	fmt.Println(duration)
}

