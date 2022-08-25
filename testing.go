package main

import (
	"fmt"
	"time"
)

func (cache *Cache) TestSet() {
	// Track speed
	startTime := time.Now()

	// Add key1 to the cache
	// Convert to map to a string before hand (json.Marshal)
	cache.Set("key1", `{
		"1":       "2",
		"summary": "my name is \"tristan\"",
	}`)
	// Print the result
	fmt.Printf("\nTest: Set() -> (%v)\n", time.Since(startTime))
}

func (cache *Cache) TestRemove() {
	// Add key2 to the map
	// Convert to map to a string before hand (json.Marshal)
	cache.Set("key2", `{
		"1":       "2",
		"summary": "my name is \"daniel\"",
	}`)
	// Track speed
	startTime := time.Now()
	// Remove key2
	cache.Remove("key2")
	// Print result
	fmt.Printf("\nTest: Remove() -> (%v)\n", time.Since(startTime))
}

func (cache *Cache) TestGet() {
	// Add key1 to the map
	// Convert to map to a string before hand (json.Marshal)
	cache.Set("key1", `{
		"1":       "2",
		"summary": "my name is \"tristan\"",
	}`)
	// Track speed
	startTime := time.Now()
	// Get the key
	var data = cache.Get("key1")
	// Print result
	fmt.Printf("\nTest: Remove() -> (%v): (%v)\n", time.Since(startTime), data)
}

// Function to test the full text search function
func (cache *Cache) TestFullTextSearch() {
	for i := 0; i < 1; i++ {
		// Convert to map to a string before hand (json.Marshal)
		cache.Set("key2", `{
			"summary": "my name is \"daniel!\"",
		}`)
	}
	// Track speed
	startTime := time.Now()

	// Get the text search result
	res := cache.FullTextSearch(TextSearch{
		Limit:      -1,
		Query:      []byte("daniel"),
		StrictMode: false,
	})
	// Print result
	fmt.Printf("FullTextSearch() Benchmark: (%v) -> %v\n\n", time.Since(startTime), res)
}
