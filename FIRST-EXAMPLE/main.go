package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup

	words := []string{ //words :=  make( []string , 4)
		                //words[0]= "alpla"
		"alpha",
		"beta",
		"gamma",
		"delta",
	}

	wg.Add(len(words))

	for i, x := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, x), &wg)
	}
	wg.Wait() //it will wait till the wait group couner becomes zero.

	wg.Add(1)
	printSomething("this is the second thing to be printed", &wg)
}
