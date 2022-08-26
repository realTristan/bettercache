package main

import "fmt"

// The Remove() function locks then unlocks the
// cache data to ensure safety before iterating through
// the cache bytes to look for the provided key
//
// once the key is found it'll search for it's closing
// bracket then remove the key from the cache bytes
//
// It will return the removed value
func (cache *Cache) Remove(key string) string {
	// Lock/Unlock the mutex
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Change the key to a modified version
	key = fmt.Sprintf(`|%s|:~`, key)

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
		if index == len(key) {

			// Set the starting index and reset the
			// key's index variable
			startIndex = i + 1
			index = 0
		} else

		// Check if the current index value is
		// equal to the key's current index
		if cache.Data[i] == key[index] {
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

				// Remove the value
				cache.Data = append(cache.Data[:startIndex-(len(key)+1)], cache.Data[i+1:]...)

				// Return the value removed
				return string(cache.Data[startIndex+2 : i])
			}
			// Reset the data length variable
			dataLength = []byte{}
		}
	}
	// Return empty string
	return ""
}
