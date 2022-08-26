package cache

import (
	"fmt"
	"strings"
	"time"
)

// The TestSet() function loops by the provided
// amount of iterations and creates a key and value
// inside the cache.
//
// Using the provided value size, it will modify
// how large the value is.
func (cache *Cache) TestSet(amountOfKeys int, valueSize int) {
	var startTime time.Time = time.Now()
	for i := 0; i < amountOfKeys; i++ {
		cache.Set(fmt.Sprintf("key%d", i), strings.Repeat("v", valueSize))
	}
	fmt.Printf("\nTest: Set(%d, %d) -> (%v)\n",
		amountOfKeys, valueSize, time.Since(startTime))
}

// The TestRemove() function tests the performance of the
// cache.Remove() function. It loops by the provided amount
// of sets to add keys to the cache
//
// After the keys are added, it then iterates over the cache
// and loops by the provided amount of removes
// and removes each key from the cache
func (cache *Cache) TestRemove(amountOfSets int, amountOfRemoves int) {
	for i := 0; i < amountOfSets; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	// Track the speed of the remove function
	var startTime time.Time = time.Now()
	for i := 0; i < amountOfRemoves; i++ {
		cache.Remove(fmt.Sprintf("key%d", i))
	}
	fmt.Printf("\nTest: Remove(%d, %d) -> (%v)\n",
		amountOfSets, amountOfRemoves, time.Since(startTime))
}

// The TestGet() function is used to test the performance of the
// Get() function. It loops by the amount of sets to add keys to
// the cache with the value size of the setValueSize parameter
//
// After the keys and values are set, it then loops over the cache
// by the amount of gets
func (cache *Cache) TestGet(amountOfSets int, setValueSize int, amountOfGets int) {
	for i := 0; i < amountOfSets; i++ {
		cache.Set(fmt.Sprintf("key%d", i), strings.Repeat("v", setValueSize))
	}
	var startTime time.Time = time.Now()
	for i := 0; i < amountOfGets; i++ {
		cache.Get(fmt.Sprintf("key%d", i))
	}
	// Print result
	fmt.Printf("\nTest: Get(%d, %d, %d) -> (%v)\n",
		amountOfSets, setValueSize, amountOfGets, time.Since(startTime))
}

// The TestFullTextSearch() function is used to test the performance
// of the cache full text search. It loops by the amount of sets to
// add keys to the cache with the value size of the setValueSize parameter
//
// After the keys and values are set, it then performs the full text search
// and whether the printResult parameter is set to true or not, it will
// print the result string.
func (cache *Cache) TestFullTextSearch(amountOfSets int, setValueSize int, printResult bool) {
	for i := 0; i < amountOfSets; i++ {
		cache.Set(fmt.Sprintf("key%d", i), strings.Repeat("v", setValueSize))
	}
	var startTime time.Time = time.Now()

	// Perform the full text search
	var result []string = cache.FullTextSearch(TextSearch{
		limit:      -1,
		query:      []byte("tristan"),
		strictMode: false,
	})
	// Print the results
	if printResult {
		fmt.Printf("Test: FullTextSearch [Result]: %v", result)
	}
	fmt.Printf("Test: FullTextSearch(%d, %d, %v) -> (%v)\n\n",
		amountOfSets, setValueSize, printResult, time.Since(startTime))
}
