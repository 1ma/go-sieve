package main

import (
	"fmt"
	"runtime"
)

var max_number, procs int

func generator(seed chan<- int) {
	for i := 2; i <= max_number; i++ {
		seed <- i
	}

	close(seed)
}

func filterer(prime int, in <-chan int, out chan<- int) {
	for number := range in {
		if number%prime != 0 {
			out <- number
		}
	}

	close(out)
}

func printer(primes <-chan int, done chan<- bool) {
	count := 0
	for prime := range primes {
		count++
		fmt.Println(prime)
	}

	fmt.Printf("\nFound %d prime numbers between 2 and %d\n", count, max_number)

	done <- true
}

func main() {
	fmt.Scanf("%d", &procs)
	fmt.Scanf("%d", &max_number)

	runtime.GOMAXPROCS(procs)

	primes := make(chan int)
	prev := make(chan int, max_number)
	done := make(chan bool)

	go printer(primes, done)
	go generator(prev)

	for {
		next := make(chan int, 1<<14)
		prime, ok := <-prev
		if !ok {
			close(primes)
			break
		}

		go filterer(prime, prev, next)

		primes <- prime
		prev = next
	}

	<-done
}
