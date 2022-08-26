package main

import (
	"fmt"
	"time"
)

// Main function
func main() {
	// Init the cache
	var Cache *Cache = Init(-1) // -1 -> Unlimited Size

	startTime := time.Now()
	// Set cache keys
	for i := 0; i < 1; i++ {
		Cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("my name is \"tristan%d\"", i))
	}

	fmt.Println(time.Since(startTime))
	// Track speed
	startTime = time.Now()

	// Get the text search result
	Cache.FullTextSearch(TextSearch{
		Limit:      -1,
		Query:      []byte("Tristan"),
		StrictMode: false,
	})
	// Print result
	fmt.Printf("Full Text Search -> (%v): %v\n\n", time.Since(startTime), nil)

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
	Cache.Show()
	fmt.Printf("Show Cache -> (%v): %v\n\n", time.Since(startTime), nil)
}
