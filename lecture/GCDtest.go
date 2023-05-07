package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Method that calculates the GCD between two numbers.
// This will need modification during the exercise
func gcd(a int, b int) int {
	fmt.Println("Calculating GCD Between ", a, " and ", b)
	i := 0
	if a > b {
		i = a
	} else {
		i = b
	}
	gcd := 0
	for gcd = i; gcd > 1; gcd-- {
		if a%gcd == 0 && b%gcd == 0 {
			return gcd
		}
	}

	return 1
}

// Right now, this code calculates the GCD of 250 pairs of numbers.
// Before starting the exercise, note how long it takes for the code
// to complete execution.

// The goal of the exercise is to parallelize this code using goroutines.
// We want to spawn three goroutines of a gcd function that all take input from a single input channel
// And send their computed GCD output to a shared output channel.

// This should be done in a way that the amount of goroutines can change without
// anything else needing modification.
// The format of how the data is printed to the console should also stay the same.

// Hint: You will want the channels to input and output information from this struct:
type gcdmsg struct {
	a   int
	b   int
	out int
}

//Example of a struct declaration: msg := gcdmsg{1,1,0}

// As expected, once the workers are correctly implemented,
// there should be a significant time saving.

// For reference, in testing with 250 pairs of numbers,
// the original solution ran in approximately 1.3 seconds on average
// While our parallelized solution with 3 workers ran in 297.23 milliseconds on average

// To run this file, type   go run GCDtest.go

func main() {
	a := 1
	b := 1
	start := time.Now()
	for k := 0; k < 250; k++ {
		a = rand.Intn(1000000)
		b = rand.Intn(1000000)
		result := gcd(a, b)
		fmt.Println("GCD Between ", a, " and ", b, " is ", result)
	}
	t := time.Now()
	duration := t.Sub(start)
	fmt.Println("Duration of the task is ", duration)

}
