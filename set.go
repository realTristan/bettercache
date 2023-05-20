package bettercache

// Import fmt Package
import "fmt"

// Set sets a new value inside the cache data.
// It locks the cache mutex to prevent data overwriting before checking whether the provided
// key already exists. If it does, it will call the Remove() function to remove that key from the cache.
//
// Once finished with the removal process, the function adds the value to the cache data slice,
// then adds the value's index to the cache indices map. If the fullText parameter is set to true,
// it also adds the value to the full text cache.
//
// Parameters:
//   - key: The cache key.
//   - value: The cache value.
//   - fullText: Whether the value should be added to the full text cache.
func (c *Cache) Set(key string, value interface{}, fullText bool) {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If key exists
	if c.existsInCache(key) || c.existsInFT(key) {
		// Remove the key from the cache
		c.remove(key)
	}

	// If the user set the AddToFullText to true
	if fullText {
		// Set the key in the cache full text indices map
		// to the index the key value is at.
		c.ft.indices[key] = len(c.ft.data)
		// Add the value into the cache data slice
		// as a modified string
		c.ft.data = append(c.ft.data, fmt.Sprintf("%v", value))
	} else {
		// Set the key in the cache indices map to the
		// index the key value is at.
		c.data[key] = value
	}
	// Increase the current cache size
	c.size++
}
