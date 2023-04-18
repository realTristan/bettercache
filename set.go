package bettercache

// Import fmt Package
import "fmt"

// The Set function is used for setting a new value inside
// the cache data. The Set function locks the cache mutex to
// prevent data overwriting before checking whether the provided
// key already exists. If it does, it will call the Remove() function
// to remove that key from the cache.
//
// Once finished with the removal process, the function
// adds the value to the cache data slice, then adds the value's
// index to the cache indices map

// Sets a key to the provided value
/* Parameters: */
/* 	key: string, { "The Cache Key" } */
/* 	value: interface{}, { "The Cache Value" } */
/* 	fullText: bool, { "Whether the value should be added to the full text cache" } */
func (c *Cache) Set(key string, value interface{}, fullText bool) {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If key exists
	if c.existsInMap(key) || c.existsInFullText(key) {
		// Remove the key from the cache
		c.remove(key)
	}

	// If the user set the AddToFullText to true
	if fullText {
		// Set the key in the cache full text indices map
		// to the index the key value is at.
		c.fullTextIndices[key] = len(c.fullTextData)
		// Add the value into the cache data slice
		// as a modified string
		c.fullTextData = append(c.fullTextData, fmt.Sprintf("%s:%v", key, value))
	} else {
		// Set the key in the cache indices map to the
		// index the key value is at.
		c.mapData[key] = value
	}
	// Increase the current cache size
	c.currentSize++
}
