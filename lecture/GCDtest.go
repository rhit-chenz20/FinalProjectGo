package main

import (
	"fmt"
	"math/rand"
	"time"
)




func gcd(a int, b int) int {
		fmt.Println("GCD Between ", a," and ", b)
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

func main() {
	a:= 1
	b:= 1
	start := time.Now()
		for k:=0; k<100;k++{
			a= rand.Intn(100000)
			b= rand.Intn(100000)
			gcd(a, b)
		}
		t := time.Now()
		duration := t.Sub(start)
		fmt.Println("Duration of the task is ", duration)

}