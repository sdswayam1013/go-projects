package main

import (
	"fmt"
	"sync"
	"time"
)

// 1-what is shared? Maps , so secure it firstby wrapping the map with a struct and guard it using mutex
type Counter struct {
	mu     sync.Mutex
	counts map[string]int
}

// Inc increments the counter for the given key.
func (c *Counter) Increment(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.counts[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *Counter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock() //defer wasnt used in increment because defer is used when somehting is returning
	return c.counts[key]
}

func main() {
	c := Counter{counts: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Increment("request")
	}

	time.Sleep(time.Second) //just to provide some extra time so that all goroutines can finish
	fmt.Println(c.Value("request"))
}
