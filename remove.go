package cache

import "fmt"

// The Remove() function locks then unlocks the
// cache data to ensure safety before iterating through
// the cache bytes to look for the provided key
//
// once the key is found it'll search for it's closing
// bracket then remove the key from the cache bytes
//
// It will return the removed value.
func (cache *Cache) Remove(key string) string {
	// Lock/Unlock the mutex
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Define Variables
	var (
		// Create a new modified key bytes
		keyBytes []byte = append([]byte(key), ':')
		// startIndex -> Track the start of the key value
		startIndex int = -1
		// index -> Track the key indexes
		index int = 0
		// Track the length of the data
		dataLength = []byte{}
	)

	// Iterate over the cache data
	for i := 1; i < len(cache.data); i++ {

		// Check if the key is present and the current
		// index is standing at that key
		if index == len(keyBytes) {

			// Set the starting index and reset the
			// key's index variable
			startIndex = i + 1
			index = 0
		} else

		// Check if the current index value is
		// equal to the key's current index
		if cache.data[i] == keyBytes[index] {
			if startIndex < 0 {
				index++
			}
		} else {
			// Reset the index
			index = 0
		}
		// Check if current index is the start of a map
		if cache.data[i] == '{' {

			// Make sure the startIndex has been established
			// and that the datalength is currently zero
			if startIndex > 0 && len(dataLength) == 0 {
				dataLength = cache.data[startIndex-1 : i]
			}
		} else

		// Check if the current index is the end of the map
		if cache.data[i] == '}' && startIndex > 0 {

			// Check if the current index is the end of the data
			if fmt.Sprint(i-startIndex-2) == string(dataLength) {
				// Remove the value
				cache.data = append(cache.data[:startIndex-(len(keyBytes)+1)], cache.data[i+1:]...)

				// Return the value removed
				return string(cache.data[startIndex+2 : i])
			}
			// Reset the data length variable
			dataLength = []byte{}
		}
	}
	// Return empty string
	return ""
}
