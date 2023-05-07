package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Method that calculates the GCD between two numbers.
// This will need modification during the exercise
func gcd(a int, b int) int {
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
			fmt.Println("GCD Between ", a," and ", b," is", gcd)
			return gcd
		}
	}
	fmt.Println("GCD Between ", a," and ", b," is 1")
	return 1
}

// Right now, this code calculates the GCD of 1000 pairs of numbers.
// Before starting the exercise, note how long it takes for the code
// to complete execution.

// The goal of the exercise is to parallelize this code using goroutines.
// We want to spawn three workers that all take input from a single input channel
// And send their computed GCD output to a shared output channel.

// This should be done in a way that the amount of workers can change without
// anything else needing modification.

// As expected, once the workers are correctly implemented, 
// there should be a significant time saving.

func main() {
	a:= 1
	b:= 1
	start := time.Now()
	for k:=0; k<250;k++{
		a = rand.Intn(1000000)
		b = rand.Intn(1000000)
		gcd(a, b)
	}
	t := time.Now()
	duration := t.Sub(start)
	fmt.Println("Duration of the task is ", duration)

}