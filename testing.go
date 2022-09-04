package bettercache

import (
	"fmt"
	"time"
)

// The TestSet() function loops by the provided
// amount of iterations and creates a key and value
// inside the cache.
//
// Using the provided value size, it will modify
// how large the value is.
func (c *Cache) TestSet(amountOfKeys int, setValue string) {
	var startTime time.Time = time.Now()
	for i := 0; i < amountOfKeys; i++ {
		c.Set(&SetData{
			Key:      fmt.Sprintf("key%d", i),
			Value:    setValue,
			FullText: true,
		})
	}
	fmt.Printf("\nTest: Set(%d, \"%s\") -> (%v)\n",
		amountOfKeys, setValue, time.Since(startTime))
}

// The TestRemove() function tests the performance of the
// c.Remove() function. It loops by the provided amount
// of sets to add keys to the cache
//
// After the keys are added, it then iterates over the cache
// and loops by the provided amount of removes
// and removes each key from the cache
func (c *Cache) TestRemove(amountOfSets int, amountOfRemoves int) {
	for i := 0; i < amountOfSets; i++ {
		c.Set(&SetData{
			Key:      fmt.Sprintf("key%d", i),
			Value:    "value",
			FullText: true,
		})
	}
	// Track the speed of the remove function
	var startTime time.Time = time.Now()
	for i := 0; i < amountOfRemoves; i++ {
		c.Remove(fmt.Sprintf("key%d", i))
	}
	fmt.Printf("\nTest: Remove(%d, %d) -> (%v)\n",
		amountOfSets, amountOfRemoves, time.Since(startTime))
}

// The TestGet() function is used to test the performance of the
// Get() function. It loops by the amount of sets to add keys to
// the cache with the provided setValue string
//
// After the keys and values are set, it then loops over the cache
// by the amount of gets
func (c *Cache) TestGet(amountOfSets int, setValue string, amountOfGets int) {
	for i := 0; i < amountOfSets; i++ {
		c.Set(&SetData{
			Key:      fmt.Sprintf("key%d", i),
			Value:    setValue,
			FullText: true,
		})
	}
	var startTime time.Time = time.Now()
	for i := 0; i < amountOfGets; i++ {
		var v interface{} = c.Get(fmt.Sprintf("key%d", i))
		fmt.Printf("%v, ", v)
	}
	// Print result
	fmt.Printf("\nTest: Get(%d, \"%s\", %d) -> (%v)\n",
		amountOfSets, setValue, amountOfGets, time.Since(startTime))
}

// The TestFullTextSearch() function is used to test the performance
// of the cache full text search. It loops by the amount of sets to
// add keys to the cache with the provided value
//
// After the keys and values are set, it then performs the full text search
// and whether the printResult parameter is set to true or not, it will
// print the result string.
func (c *Cache) TestFullTextSearch(
	amountOfSets int, setValue string, searchFor string, printResult bool) {

	// Set cache data
	for i := 0; i < amountOfSets; i++ {
		c.Set(&SetData{
			Key:      fmt.Sprintf("key%d", i),
			Value:    setValue,
			FullText: true,
		})
	}
	// Track speed
	var startTime time.Time = time.Now()

	// Perform the full text search.
	var result []string = c.FullTextSearch(&TextSearch{
		Limit:      -1,
		Query:      searchFor,
		StrictMode: false,
	})
	var endTime time.Duration = time.Since(startTime)

	// Print the results
	if printResult {
		fmt.Printf("Test: FullTextSearch [Result]: %v\n", result)
	}
	fmt.Printf("Test: FullTextSearch(%d, %s, %s, %v) -> (%v)\n\n",
		amountOfSets, setValue, searchFor, printResult, endTime)
}
