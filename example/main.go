package main

import (
	"fmt"
	"time"
)

func main() {
	cache := NewCache(5 * time.Second)

	cache.Set("foo", "bar")
	val, ok := cache.Get("foo")
	fmt.Println(val, ok) // bar true

	time.Sleep(6 * time.Second)
	val, ok = cache.Get("foo")
	fmt.Println(val, ok) // <nil> false (expired)
}
