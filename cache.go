package bettercache

// Import Packages
import (
	"fmt"
	"sync"
)

// FullText is a struct that holds the data and indices for full text search.
// It contains two fields:
//   - data: a slice of strings that holds the cache values for full text search.
//   - indices: a map that holds the indices of the cache values in the data slice.
type FullText struct {
	data    []string
	indices map[string]int
}

// The Cache struct is used for storing and managing cache data.
// It contains four fields:
//   - size: an integer that holds the current size of the cache.
//   - mutex: a pointer to a sync.RWMutex that is used for locking and unlocking the cache data.
//   - ft: a pointer to a FullText struct that holds the data and indices for full text search.
//   - data: a map that holds the main data cache values.
type Cache struct {
	size  int
	mutex *sync.RWMutex
	ft    *FullText
	data  map[string]interface{}
}

// The InitCache function is used for initializing a new cache instance.
// It returns a pointer to a new Cache struct with its fields initialized.
//
// Returns:
//   - A pointer to a new Cache struct with its fields initialized.
func InitCache() *Cache {
	return &Cache{
		mutex: &sync.RWMutex{},
		size:  0,
		data:  make(map[string]interface{}),
		ft: &FullText{
			data:    make([]string, 0),
			indices: make(map[string]int),
		},
	}
}

// The Exists function is used for checking whether a key
// exists in the cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked.
//
// Parameters:
//   - key: The cache key to check for existence.
//
// Returns:
//   - A boolean value indicating whether the provided key exists in the cache.
func (c *Cache) Exists(key string) bool {
	// Mutex locking/unlocking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Checks the full text cache and the data map
	return c.existsInFT(key) || c.existsInCache(key)
}

// The Get function is used for retrieving a cache value
// using its corresponding key. The function read locks the
// cache mutex before checking whether the key exists in the
// cache. If the key exists in the full text data, the function
// returns the full text cache string. If the key exists in the
// data map, the function returns the map data value. Once the
// function returns, the mutex is unlocked.
//
// Parameters:
//   - key: The cache key to retrieve the value for.
//
// Returns:
//   - The cache value of the provided key, or nil if the key does not exist in the cache.
func (c *Cache) Get(key string) interface{} {
	// Mutex locking/unlocking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Check if the key exists in the full text data
	if c.existsInFT(key) {
		// Return the full text cache string
		return c.ft.data[c.ft.indices[key]][len(fmt.Sprint(key))+1:]
	} else

	// The key exists in the cache data map
	if c.existsInCache(key) {
		// Return the map data value
		return c.data[key]
	}
	// Return nil
	return nil
}

// The Remove function is used for removing a cache value
// using its corresponding key. The function write locks the
// cache mutex before removing the key-value pair from the
// cache data map. Once the key-value pair is removed, the
// function moves onto iterating through the c.ft.indices map
// reducing all cache indices that are post-the-removed-value.
// The function then returns the removed value. Once the function
// returns, the cache mutex is unlocked.
//
// Parameters:
//   - key: The cache key to remove the value for.
//
// Returns:
//   - The removed cache value of the provided key, or nil if the key does not exist in the cache.
func (c *Cache) Remove(key string) interface{} {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Return the removed key
	return c.remove(key)
}

// The Walk function is used for iterating over the cache data.
// The function read locks the cache mutex before iterating over
// the data map or full text indices. If the full text flag is true,
// the function iterates over the full text indices and calls the
// provided function with the key and full text value as arguments.
// If the full text flag is false, the function iterates over the
// data map and calls the provided function with the key and value
// as arguments. If the provided function returns false, the iteration
// stops. Once the iteration is complete, the mutex is unlocked.
//
// Parameters:
//   - ft: A boolean flag indicating whether to iterate over the full text data or not.
//   - fn: A function that takes a key and value as arguments and returns a boolean value.
//     The function is called for each key-value pair in the data map or full text indices.
//
// Example:
//
//	cache.Walk(false, func(k string, v interface{}) bool {
//	  fmt.Printf("%s: %v\n", k, v)
//	  return true
//	})
func (c *Cache) Walk(ft bool, fn func(key string, val interface{}) bool) {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Check if the full text flag is true
	if ft {
		// Iterate over the full text indices
		for k, v := range c.ft.indices {
			// get the full text data using
			// the full text index
			var _v string = c.ft.data[v]
			// Call the provided function
			if !fn(k, _v) {
				return
			}
		}
	} else {
		for k, v := range c.data {
			if !fn(k, v) {
				return
			}
		}
	}
}

// The Keys function is used for returning a slice of all the keys in the cache.
// The function read locks the cache mutex before iterating over the data map.
// For each key in the data map, the function appends the key to a slice of keys.
// Once the iteration is complete, the mutex is unlocked and the slice of keys is returned.
//
// Returns:
//   - A slice of all the keys in the cache.
func (c *Cache) Keys() []string {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Define Variables
	var (
		// keys -> The slice containing the keys
		keys []string = []string{}
		// i -> Track the index for setting the keys
		i int = 0
	)
	// Iterate over the cache indices map
	for k := range c.data {
		keys[i] = k
		i++
	}
	// Return the keys slice
	return keys
}

// The Clear function is used to clear the cache data and the cache indices.
// The function locks the cache mutex before resetting the cache data and the
// cache indices. Once the function returns, the mutex is unlocked.
//
// Example:
//
//	cache.Clear()
func (c *Cache) Clear() {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()
	*c = *InitCache()
}

// The Size function is used for getting the current size of the cache.
// The function read locks the cache mutex before returning the cache size.
// Once the function returns, the cache mutex is unlocked.
//
// Returns:
//   - The current size of the cache as an integer.
func (c *Cache) Size() int {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Create a copy of the size
	var copy int = c.size

	// Return the copy
	return copy
}

// The existsInFT function is used for checking whether a key exists in the full text cache or not.
// The function read locks the cache mutex before returning whether the key is in the cache.
// Once the function returns, the mutex is unlocked.
//
// Parameters:
//   - key: A string representing the cache key to check.
//
// Returns:
//   - A boolean value indicating whether the provided key exists in the full text cache or not.
func (c *Cache) existsInFT(key string) bool {
	// Check if the key exists within the c.ft.indices
	if _, t := c.ft.indices[key]; t {
		return true
	}
	return false
}

// The existsInCache function is used for checking whether a key exists in the main cache or not.
// The function read locks the cache mutex before returning whether the key is in the cache.
// Once the function returns, the mutex is unlocked.
//
// Parameters:
//   - key: A string representing the cache key to check.
//
// Returns:
//   - A boolean value indicating whether the provided key exists in the main cache or not.
func (c *Cache) existsInCache(key string) bool {
	// Check if the key exists within the c.ft.indices
	if _, t := c.data[key]; t {
		return true
	}
	return false
}

// The remove function is used for removing a key-value pair from the cache.
// The function checks whether the key exists in the main cache or the full text cache.
// If the key exists in the main cache, the function removes the key-value pair from the main cache.
// If the key exists in the full text cache, the function removes the key-value pair from the full text cache.
// The function returns the removed value if the key exists in the cache, otherwise it returns nil.
//
// Parameters:
//   - key: A string representing the cache key to remove.
//
// Returns:
//   - The removed value if the key exists in the cache, otherwise it returns nil.
func (c *Cache) remove(key string) interface{} {
	// Check if the value is not a full text value
	if c.existsInCache(key) {
		c.size--

		// Create the remvoedValue for return
		var v interface{} = c.data[key]

		// Delete key from c.mainIndices map
		delete(c.data, key)

		// Make sure the cache data isn't empty
		if len(c.data) > 0 {
			// Return the removed value
			return v
		}
	} else if c.existsInFT(key) {
		c.size--

		// Remove key from the cache slice
		c.ft.data = append(c.ft.data[:c.ft.indices[key]],
			c.ft.data[c.ft.indices[key]+1:]...)

		// Iterate over the map keys
		for k := range c.ft.indices {
			if c.ft.indices[k] > c.ft.indices[key] {
				c.ft.indices[k] -= 1
			}
		}
		// Delete key from c.ft.indices map
		delete(c.ft.indices, key)

		// Make sure the cache data isn't empty
		if len(c.ft.data) > 0 {
			// Return the removed value
			return c.ft.data[c.ft.indices[key]]
		}
	}
	// Return nil
	return nil
}
