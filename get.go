package cache

import (
	"bytes"
	"fmt"
)

// The Get() function read locks then read unlocks
// the cache data to ensure safety before returning
// a json map with the key's value.
func (cache *Cache) Get(key string) string {
	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Make sure the cache contains the provided key first
	var keyBytes []byte = append([]byte(key), ':')
	if !bytes.Contains(cache.Data, keyBytes) {
		return ""
	}

	// Define Variables
	var (
		// startIndex -> Track the start of the key value
		startIndex int = -1
		// index -> Track the key indexes
		index int = 0
		// Track the length of the data
		dataLength = []byte{}
	)

	// Iterate over the cache data
	for i := 1; i < len(cache.Data); i++ {
		// Check if the key is present and the current
		// index is standing at that key
		if index == len(keyBytes) {

			// Set the starting index and reset the
			// key's index variable
			startIndex = i + 1
			index = 0
		} else

		// Check if the current index value is
		// equal to the keyBytes's current index
		if cache.Data[i] == keyBytes[index] {
			// Make sure the startIndex has been established
			if startIndex < 0 {
				index++
			}
		} else {
			// Reset the index
			index = 0
		}

		// Check if current index is the start of a map
		if cache.Data[i] == '{' {

			// Make sure the startIndex has been established
			// and that the datalength is currently zero
			if startIndex > 0 && len(dataLength) == 0 {
				dataLength = cache.Data[startIndex-1 : i]
			}
		} else

		// Check if the current index is the end of the map
		if cache.Data[i] == '}' && startIndex > 0 {

			// Check if the current index is the end of the data
			if fmt.Sprint(i-startIndex-2) == string(dataLength) {
				// Return the data
				return string(cache.Data[startIndex+2 : i])
			}
			// Reset the data length variable
			dataLength = []byte{}
			startIndex = -1
		}
	}
	// Return empty string
	return ""
}
