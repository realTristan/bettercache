package bettercache

// Import Packages
import (
	"fmt"
	"sync"
)

// The Cache struct has six primary keys (all unexported)
/* currentSize: int { "The current map size" } */
/* mutex: *sync.RWMutex { "The mutex for locking/unlocking the data" } 				  */
/* mapData: map[interface{}]interface{} { "The Main Data Cache Values" } 								  */
/* fullTextData: []string { "The Full Text Data Cache Values" } 					  */
/* fulltextIndices: map[string]int { "The Cache Keys holding the full text indices of the Cache Values" } 	*/
type Cache struct {
	currentSize     int
	mutex           *sync.RWMutex
	fullTextData    []string
	fullTextIndices map[interface{}]int
	mapData         map[interface{}]interface{}
}

// The Exists function is used for checking whether a key
// exists in the cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	doesExist: bool 						*/
func (c *Cache) Exists(key interface{}) bool {
	// Mutex locking/unlocking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Checks the full text cache and the data map
	return c.existsInFullText(key) || c.existsInMap(key)
}

// The Get function is used for return a value from the cache
// with a key. The function read locks the cache mutex before
// checking whether the key exists in the cache. If the key
// does exist, it will use the cache indices map to get the cache data
// index and return the cache value. Once the function returns,
// the mutex is unlocked
//
// If the cache value's FullText has been set to true, it will split
// the value by ':' and return the index[2] of it's result

// Returns the cache value of the provided key
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	cacheValue: interface{} 				*/
func (c *Cache) Get(key interface{}) interface{} {
	// Mutex locking/unlocking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Check if the key exists in the full text data
	if c.existsInFullText(key) {
		// Return the full text cache string
		return c.fullTextData[c.fullTextIndices[key]][len(fmt.Sprint(key))+1:]
	} else

	// The key exists in the cache data map
	if c.existsInMap(key) {
		// Return the map data value
		return c.mapData[key]
	}
	// Return nil
	return nil
}

// The remove() function is used to remove a value from the cache
// using it's corresponding key. The function full locks the mutex
// before modifying the c.mapData, removing the cache value.
//
// Once the c.mapData value is removed, the function moves onto
// iterating through the c.fullTextIndices map reducing all cache indices
// that are post-the-removed-value. The function then returns the
// removed value. Once the function returns, the cache mutex is unlocked.

// Removes a key from the cache
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	removedValue: interface{} 				*/
func (c *Cache) Remove(key interface{}) interface{} {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Return the removed key
	return c.remove(key)
}

// The Init function is used for initalizing the Cache struct
// and returning the newly created cache object.
// You can create the cache by yourself, but using the Init()
// function is much easier

// Initializes the cache object
//
/* Returns 													*/
/* 	cache: *Cache 											*/
func InitCache() *Cache {
	return initCache()
}

// The WalkNonFT function is used for iterating over
// the non-full text search cache data.
/*
 Example:
	cache.Walk(func(k, v string) bool {
		fmt.Printf("%s: %s\n", k, v)
		return true
	}
*/
func (c *Cache) WalkNonFT(fn func(key, val interface{}) bool) {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// You can use a lock for each "get" here
	for k, v := range c.mapData {
		if !fn(k, v) {
			return
		}
	}
}

// The WalkFT function is used for iterating over
// the full text search cache data.
/*
 Example:
	cache.Walk(func(k, v string) bool {
		fmt.Printf("%s: %s\n", k, v)
		return true
	}
*/
func (c *Cache) WalkFT(fn func(key, val interface{}) bool) {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Iterate over the full text indices
	for k, v := range c.fullTextIndices {
		// get the full text data using
		// the full text index
		var _v string = c.fullTextData[v]
		// Call the provided function
		if !fn(k, _v) {
			return
		}
	}
}

// The ShowKeys function is used to get a slice of all
// the cache keys. The function read locks the cache mutex
// before iterating over the cache indices map, adding each
// of the keys to the keys slice.
//
// The function then returns the slice of keys. Once the
// function returns, the cache mutex is unlocked

// Returns all the cache keys in a slice
//
/* Returns 							*/
/* 	keys: []interface{}				*/
func (c *Cache) ShowKeys() []interface{} {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Define Variables
	var (
		// keys -> The slice containing the keys
		keys []interface{} = make([]interface{}, len(c.mapData))
		// i -> Track the index for setting the keys
		i int = 0
	)
	// Iterate over the cache indices map
	for k := range c.mapData {
		keys[i] = k
		i++
	}
	// Return the keys slice
	return keys
}

// The Clear function is used to clear the cache
// data and the cache indices. The function locks
// the cache mutex before resetting the cache data
// and the cache indices. Once the function
// returns, the mutex is unlocked

// Clear the cache data
func (c *Cache) Clear() {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()

	*c = *(initCache())
}

// The CurrentSize function is used to get the current size
// of the cache. The function read locks the cache mutex
// before returning the cache currentSize. Once the function
// returns, the cache mutex is unlocked.

// Return the cache current size (int)
func (c *Cache) CurrentSize() int {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Copy the caches max size using a function
	// then return it. This is to accompany
	// for safety.
	return func(cs int) int { return cs }(c.currentSize)
}

// The initCache function is used for initalizing the Cache struct
// and returning the newly created cache object.
// You can create the cache by yourself, but using the Init()
// function is much easier

// Initializes the cache object
//
/* Returns 													*/
/* 	cache: *Cache 											*/
func initCache() *Cache {
	// Return a new Cache object
	return &Cache{
		mutex:           &sync.RWMutex{},
		currentSize:     0,
		mapData:         make(map[interface{}]interface{}),
		fullTextData:    []string{},
		fullTextIndices: make(map[interface{}]int),
	}
}

// The existsinFullText function is used for checking whether a key
// exists in the full text cache or not. The function read locks
// the cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	doesExist: bool 						*/
func (c *Cache) existsInFullText(key interface{}) bool {
	// Check if the key exists within the c.fullTextIndices
	if _, t := c.fullTextIndices[key]; t {
		return true
	}
	return false
}

// The existsInMap function is used for checking whether a key
// exists in the main cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
/* Parameters: 							*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 								*/
/* 	doesExist: bool 					*/
func (c *Cache) existsInMap(key interface{}) bool {
	// Check if the key exists within the c.fullTextIndices
	if _, t := c.mapData[key]; t {
		return true
	}
	return false
}

// The remove() function is used to remove a value from the cache
// using it's corresponding key. The function full locks the mutex
// before modifying the c.mapData, removing the cache value.
//
// Once the c.mapData value is removed, the function moves onto
// iterating through the c.fullTextIndices slice reducing all cache indices
// that are post-the-removed-value. The function then returns the
// removed value. Once the function returns, the cache mutex is unlocked.

// Removes a key from the cache
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	removedValue: interface{} 				*/
func (c *Cache) remove(key interface{}) interface{} {
	// Check if the value is not a full text value
	if c.existsInMap(key) {
		c.currentSize--

		// Create the remvoedValue for return
		var removedValue interface{} = c.mapData[key]

		// Delete key from c.mainIndices map
		delete(c.mapData, key)

		// Make sure the cache data isn't empty
		if len(c.mapData) > 0 {
			// Return the removed value
			return removedValue
		}
	} else

	// Else if the key exists in the full text data
	if c.existsInFullText(key) {
		c.currentSize--

		// Remove key from the cache slice
		c.fullTextData = append(c.fullTextData[:c.fullTextIndices[key]],
			c.fullTextData[c.fullTextIndices[key]+1:]...)

		// Iterate over the map keys
		for k := range c.fullTextIndices {
			if c.fullTextIndices[k] > c.fullTextIndices[key] {
				c.fullTextIndices[k] -= 1
			}
		}
		// Delete key from c.fullTextIndices map
		delete(c.fullTextIndices, key)

		// Make sure the cache data isn't empty
		if len(c.fullTextData) > 0 {
			// Return the removed value
			return c.fullTextData[c.fullTextIndices[key]]
		}
	}
	// Return nil
	return nil
}
