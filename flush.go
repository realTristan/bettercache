package cache

import (
	"fmt"
	"os"
)

func getWrapper(o interface{}) string {
	switch o.(type) {
	case int:
		return fmt.Sprintf("$INT(%v", o)
	}
	return fmt.Sprintf("$STR(%v", o)
}

// The flushToFile function is used to write all the
// cache data to the BetterCache file.
// In the future, there will be more options for
// this file but as for now it will write to the file
// whether the cache key+value is a full text key or a
// regular cache map key. As well as it will write
// the key and value
//
/* Paramters: */
/* path: string { "The path to the BetterCache file" } */
func (cache *Cache) flushToFile(path string) {
	// If the BetterCache file doesn't exist
	if _, err := os.Stat(path); err != nil {
		// Return the function
		return
	}

	// Create the result map
	var result []byte = []byte{}

	// Make sure the full text indices is greater than
	// zero. Having a check prevents potential bugs.
	if len(cache.fullTextIndices) > 0 {
		// Iterate over the full text indices
		for _, i := range cache.fullTextIndices {

			// Append the full text key and value to the
			// result byte array
			result = append(result,
				[]byte(fmt.Sprintf("$FULLTEXT:%s\n", cache.fullTextData[i]))...)
		}
	}

	// Make sure the cache map data is greater than
	// zero. Having a check prevents potential bugs.
	if len(cache.mapData) > 0 {

		// Iterate over the cache map data
		for k, v := range cache.mapData {
			// Append the cache map key and value to the
			// result byte array
			result = append(result,
				[]byte(fmt.Sprintf("$DATA:%v):%s)\n", getWrapper(k), getWrapper(v)))...)
		}
	}

	// Write the result byte array to the BetterCache file
	os.WriteFile(path, result, 0644)
}

// The Flush function locks the mutex before calling
// the flushToFile function. Once the function has been
// called and the Flush function returns, the mutex is unlocked

// The Flush function is used to write the cache
// data to a BetterCache file
//
/* Paramters: */
/* 	path: string { "The path to the BetterCache file" } */
func (cache *Cache) Flush(path string) {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Flush cache to BetterCache file
	cache.flushToFile(path)
}
