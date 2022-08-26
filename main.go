package main

import (
	"fmt"
	"time"
)

// Main function
func main() {
	// Init the cache
	var Cache *Cache = Init(-1) // -1 -> Unlimited Size

	// Set cache keys
	for i := 0; i < 2; i++ {
		Cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("my name is \"tristan%d\"", i))
	}
	// Track speed
	startTime := time.Now()

	// Get the text search result
	res := Cache.FullTextSearch(TextSearch{
		Limit:      -1,
		Query:      []byte("tristan"),
		StrictMode: false,
	})
	// Print result
	fmt.Printf("Full Text Search -> (%v): %v\n\n", time.Since(startTime), res)

	// Test setting a duplicate key
	startTime = time.Now()
	Cache.Set("key1", "my name is \"daniel\"")
	fmt.Printf("Set (key1) -> (%v)\n\n", time.Since(startTime))

	// Test getting a key
	startTime = time.Now()
	k := Cache.Get("key1")
	fmt.Printf("Get Key -> (%v): %s\n\n", time.Since(startTime), k)

	// Test showing the cache
	startTime = time.Now()
	s := Cache.Show()
	fmt.Printf("Show Cache -> (%v): %v\n\n", time.Since(startTime), s)
}
