package cache

// Import Packages
import (
	"fmt"
	"sync"
)

// The _Cache struct has six primary keys (all unexported)
/* currentSize: int { "The current map size" } */
/* maxSize: int { "The maximum map size" } */
/* mutex: *sync.RWMutex { "The mutex for locking/unlocking the data" } 				  */
/* mapData: map[interface{}]interface{} { "The Main Data Cache Values" } 								  */
/* fullTextData: []string { "The Full Text Data Cache Values" } 					  */
/* fulltextIndices: map[string]int { "The Cache Keys holding the full text indices of the Cache Values" } 	*/
type _Cache struct {
	currentSize     int
	maxSize         int
	mutex           *sync.RWMutex
	fullTextData    []string
	fullTextIndices map[interface{}]int
	mapData         map[interface{}]interface{}
}

// The Init function is used for initalizing the Cache struct
// and returning the newly created cache object.
// You can create the cache by yourself, but using the Init()
// function is much easier

// Initializes the cache object
/* Parameters: 												*/
/* 	size: int { "The Size of the cache map and slice" }  	*/
//
/* Returns 													*/
/* 	cache: *_Cache 											*/
func Init(size int) *_Cache {
	var cache *_Cache = &_Cache{
		mutex:       &sync.RWMutex{},
		maxSize:     size,
		currentSize: 0,
	}
	// If the user passed a size variable greater
	// tham zero.
	if size > 0 {
		// If so, it will use the make() function
		// for creating a data slice and index map
		// with a set size
		cache.mapData = make(map[interface{}]interface{}, size)
		cache.fullTextData = make([]string, size)
		cache.fullTextIndices = make(map[interface{}]int, size)
	}
	// Return a cache with no limit to it's
	// size. It is recommended that you provide a size.
	cache.mapData = map[interface{}]interface{}{}
	cache.fullTextData = []string{}
	cache.fullTextIndices = make(map[interface{}]int)

	// Return the cache
	return cache
}

// The ExistsinFullText function is used for checking whether a key
// exists in the full text cache or not. The function read locks
// the cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	doesExist: bool 						*/
func (cache *_Cache) ExistsInFullText(key interface{}) bool {
	// Check if the key exists within the cache.fullTextIndices
	if _, t := cache.fullTextIndices[key]; t {
		return true
	}
	return false
}

// The ExistsInMap function is used for checking whether a key
// exists in the main cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
/* Parameters: 							*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 								*/
/* 	doesExist: bool 					*/
func (cache *_Cache) ExistsInMap(key interface{}) bool {
	// Check if the key exists within the cache.fullTextIndices
	if _, t := cache.mapData[key]; t {
		return true
	}
	return false
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
func (cache *_Cache) Exists(key interface{}) bool {
	// Checks the full text cache and the data map
	return cache.ExistsInFullText(key) || cache.ExistsInMap(key)
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
func (cache *_Cache) Get(key interface{}) interface{} {
	// Check if the key exists in the full text data
	if cache.ExistsInFullText(key) {
		// Mutex locking/unlocking
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		// Return the full text cache string
		return cache.fullTextData[cache.fullTextIndices[key]][len(fmt.Sprint(key))+1:]
	} else

	// The key exists in the cache data map
	if cache.ExistsInMap(key) {
		// Mutex locking/unlocking
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		// Return the map data value
		return cache.mapData[key]
	}
	// Return nil
	return nil
}

// The Remove function is used to remove a value from the cache
// using it's corresponding key. The function full locks the mutex
// before modifying the cache.mainData, removing the cache value.
//
// Once the cache.mapData value is removed, the function moves onto
// iterating through the cache.mainIndices map reducing all cache indices
// that are post-the-removed-value. The function then returns the
// removed value. Once the function returns, the cache mutex is unlocked.

// Removes a key from the cache
/* Parameters: 								*/
/* 	key: interface{} { "The Cache Key" } 	*/
//
/* Returns 									*/
/* 	removedValue: interface{} 				*/
func (cache *_Cache) Remove(key interface{}) interface{} {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Check if the value is not a full text value
	if cache.ExistsInMap(key) {
		cache.currentSize--

		// Create the remvoedValue for return
		var removedValue interface{} = cache.mapData[key]

		// Delete key from cache.mainIndices map
		delete(cache.mapData, key)

		// Make sure the cache data isn't empty
		if len(cache.mapData) > 0 {
			// Return the removed value
			return removedValue
		}
	} else

	// Else if the key exists in the full text data
	if cache.ExistsInFullText(key) {
		cache.currentSize--

		// Remove key from the cache slice
		cache.fullTextData = append(cache.fullTextData[:cache.fullTextIndices[key]],
			cache.fullTextData[cache.fullTextIndices[key]+1:]...)

		// Iterate over the map keys
		for k := range cache.fullTextIndices {
			if cache.fullTextIndices[k] > cache.fullTextIndices[key] {
				cache.fullTextIndices[k] -= 1
			}
		}
		// Delete key from cache.fullTextIndices map
		delete(cache.fullTextIndices, key)

		// Make sure the cache data isn't empty
		if len(cache.fullTextData) > 0 {
			// Return the removed value
			return cache.fullTextData[cache.fullTextIndices[key]]
		}
	}
	// Return nil
	return nil
}

// The Show function is used for getting the cache data.
// The function read locks the mutex then returns the
// cache.mapData and the cache.fullTextData,
// Once the function has returned, the mutex unlocks

// Show the cache
/* Returns 								*/
/* 	cache.mainData: []interface{} 		*/
func (cache *_Cache) Show() (map[interface{}]interface{}, []string) {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the cache values
	return cache.mapData, cache.fullTextData
}

// The ShowFTIndices function is used for getting the cache
// slice indices. The function read locks the mutex
// then returns the cache.fullTextIndices. Once the function has returned,
// the mutex unlocks

// Show the cache
/* Returns 								*/
/* 	cache.mainIndices: map[interface{}]int 		*/
func (cache *_Cache) ShowFTIndices() map[interface{}]int {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the cache indices map
	return cache.fullTextIndices
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
func (cache *_Cache) ShowKeys() []interface{} {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Define Variables
	var (
		// keys -> The slice containing the keys
		keys []interface{} = make([]interface{}, len(cache.mapData))
		// i -> Track the index for setting the keys
		i int = 0
	)
	// Iterate over the cache indices map
	for k := range cache.mapData {
		keys[i] = k
		i++
	}
	// Return the keys slice
	return keys
}

// The Clear function is used to clear the cache
// data and the cache indices. The function locks
// the cache mutex before resetting the cache data
// and the cache indices. Once the function returns
// the cache mutex is unlocked

// Clear the cache data
func (cache *_Cache) Clear() *_Cache {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Reset the cache variables
	return Init(cache.maxSize)
}

// The GetMaxSize function is used to get the maximum size
// of the cache. The function read locks the cache mutex
// before returning the cache maxSize. Once the function
// returns, the cache mutex is unlocked.

// Returns the caches maximum size (int)
func (cache *_Cache) GetMaxSize() int {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the caches maximum size
	return cache.maxSize
}

// The GetCurrentSize function is used to get the current size
// of the cache. The function read locks the cache mutex
// before returning the cache currentSize. Once the function
// returns, the cache mutex is unlocked.

// Return the cache current size (int)
func (cache *_Cache) GetCurrentSize() int {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the caches maximum size
	return cache.currentSize
}
