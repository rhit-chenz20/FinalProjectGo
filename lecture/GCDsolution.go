package main

import (
	"fmt"
	"math/rand"
	"time"
)

type gcdmsg struct {
	a	int
	b	int
	out	int
}

func gcd(cin chan gcdmsg, cout chan gcdmsg) {
	for {
		msg := <- cin
		a := msg.a
		b := msg.b
		fmt.Println("Calculating GCD Between ", a," and ", b)
		i:=0
		if a>b {
			i = a
		}else{
			i=b
		}
		gcd:=0
		for gcd=i;gcd>1;gcd--{
			if a%gcd == 0 && b%gcd == 0{
				msg.out = gcd
				cout <- msg
			}
		}
		msg.out = 1
		cout <- msg
	}
}

const TEST_LEN = 1000
const NUM_WORKERS = 10

func main() {
	msg := gcdmsg{1,1,0}

	cin := make(chan gcdmsg, TEST_LEN)
	cout := make(chan gcdmsg, TEST_LEN)

	// Create the workers
	for i := 0; i < NUM_WORKERS; i++ {
		go gcd(cin, cout)
	}

	start := time.Now()
	// Fill up the Input Queue
	for i := 0; i < TEST_LEN; i++ {
		msg.a = rand.Intn(1000000)
		msg.b = rand.Intn(1000000)
		cin <- msg
	}

	// Read the answers from the Output channel
	for i := 0; i < TEST_LEN; i++ {
		val := <-cout
		fmt.Println("GCD Between ", val.a," and ", val.b," is ", val.out)
	}

	t := time.Now()
	duration := t.Sub(start)
	fmt.Println("Duration of the task is ", duration)
}