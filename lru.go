package main

import (
	"fmt"
	"math/rand"
)

func main() {
	cache := new(LRUCache)
	cache.init(5)
	for i := 0; i < 100; i++ {
		num := rand.Int() % 20
		fmt.Printf("Writing: %d, ", num)
		evicted, _ := cache.write(fmt.Sprint(num), num)
		fmt.Printf("evicted: %v, latest:%+v, size: %d\n", evicted, cache.Q.tail.data, cache.size)
		fmt.Printf("%s\n-----------------\n", cache.Q.String())
	}
}
