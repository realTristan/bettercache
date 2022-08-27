package cache

// Import Packages
import (
	"fmt"
	"sync"
)

// The Cache struct has three primary keys (all unexported)
/* data: []string { "The Cache Values" } 											  */
/* mutex: *sync.RWMutex { "The mutex for locking/unlocking the data" } 				  */
/* indices: map[string]int { "The Cache Keys holding the indices of the Cache Values" } */
type Cache struct {
	mutex   *sync.RWMutex
	data    []interface{}
	indices map[string]int
}

// The SetData struct has three primary keys
/* Key: string { "The Cache Key" }		*/
/* Value: string { "The Cache Value" }	*/
/* FullText: string { "Whether to enable full text functions " }	*/
// WARNING
/* If FullText is set to true, it converts the Value to a string	*/
type SetData struct {
	Key      string
	Value    interface{}
	FullText bool
}

// The Init function is used for initalizing the Cache struct
// and returning the newly created cache object.
// You can create the cache by yourself, but using the Init()
// function is much easier

// Initializes the cache object
/* Parameters 												*/
/* 	size: int { "The Size of the cache map and slice" }  	*/
//
/* Returns 													*/
/* 	cache: *Cache 											*/
func Init(size int) *Cache {
	// If the user passed a size variable greater
	// tham zero.
	if size > 0 {
		// If so, it will use the make() function
		// for creating a data slice and index map
		// with a set size
		return &Cache{
			mutex:   &sync.RWMutex{},
			data:    make([]interface{}, size),
			indices: make(map[string]int, size),
		}
	}
	// Return a cache with no limit to it's
	// size. It is recommended that you provide a size.
	return &Cache{
		mutex:   &sync.RWMutex{},
		data:    []interface{}{},
		indices: make(map[string]int),
	}
}

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
/* Parameters 										 			    		*/
/* 	key: string { "The Cache Key" }  					 			    	*/
/* 	value: interface{} { "The value to set the key to" }  			    	*/
/* 	fullText: bool { "Whether to add the value to the full text slice" } 	*/
func (cache *Cache) Set(SD *SetData) {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// If key exists
	if _, t := cache.indices[SD.Key]; t {
		// I decided to put this inside a function so that
		// even if there's any errors in the Remove function,
		// the mutex will still relock once the function returns
		func() {
			// Unlock the mutex so the remove function can
			// remove the key from the cache
			cache.mutex.Unlock()
			// Then re-lock the mutex once the key has
			// been removed
			defer cache.mutex.Lock()

			// Remove the key from the cache
			cache.Remove(SD.Key)
		}()
	}

	// Set the key in the cache indices map to the
	// index the key value is at.
	cache.indices[SD.Key] = len(cache.data)

	// If the user set the AddToFullText to true
	if SD.FullText {
		// Add the value into the cache data slice
		// as a modified string
		cache.data = append(cache.data,
			fmt.Sprintf("FT(true):%s:%v", SD.Key, SD.Value))
	} else {
		// Else
		// Add the value into the cache data slice
		cache.data = append(cache.data, SD.Value)
	}
}

// The Exists function is used for checking whether a key
// exists in the cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
/* Parameters 							*/
/* 	key: string { "The Cache Key" } 	*/
//
/* Returns 								*/
/* 	doesExist: bool 					*/
func (cache *Cache) Exists(key string) bool {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Check if the key exists within the cache.indices
	if _, t := cache.indices[key]; t {
		// It does
		return true
	}
	// It does not
	return false
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
/* Parameters 						*/
/* 	key: string { "The Cache Key" } */
//
/* Returns 							*/
/* 	cacheValue: interface{} 		*/
func (cache *Cache) Get(key string) interface{} {
	// Make sure the key exists before returning
	// the key's value
	// If you don't check whether the key exists before
	// then it will return the key prior to the given in
	// the cache.indices map
	if cache.Exists(key) {
		// Mutex locking
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		// Return the key value
		//if strings.StartsWith("FT(true)") {
		//	return strings.Split(cache.data[cache.indices[key]].(string), ":")[2]
		//}
		return cache.data[cache.indices[key]]
	}
	// Return empty string
	return ""
}

// The Get function is used for return a value from the cache
// with a key. The function read locks the cache mutex before
// checking whether the key exists in the cache. If the key
// does exist, it will use the cache indices map to get the cache data
// index and return the cache value. Once the function returns,
// the mutex is unlocked

// Returns the unmodified cache value of the provided key
/* Parameters 						*/
/* 	key: string { "The Cache Key" } */
//
/* Returns 							*/
/* 	cacheValue: interface{} 		*/
func (cache *Cache) GetFull(key string) interface{} {
	// Make sure the key exists before returning
	// the key's value
	// If you don't check whether the key exists before
	// then it will return the key prior to the given in
	// the cache.indices map
	if cache.Exists(key) {
		// Mutex locking
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		// Return the key value
		return cache.data[cache.indices[key]]
	}
	// Return empty string
	return ""
}

// The Remove function is used to remove a value from the cache
// using it's corresponding key. The function full locks the mutex
// before modifying the cache.data, removing the cache value.
//
// Once the cache.data value is removed, the function moves onto
// iterating through the cache.indices map reducing all cache indices
// that are post-the-removed-value. The function then returns the
// removed value. Once the function returns, the cache mutex is unlocked.

// Removes a key from the cache
/* Parameters 						*/
/* 	key: string { "The Cache Key" } */
//
/* Returns 							*/
/* 	removedValue: interface{} 		*/
func (cache *Cache) Remove(key string) interface{} {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Remove key from the cache slice
	cache.data = append(cache.data[:cache.indices[key]],
		cache.data[cache.indices[key]+1:]...)

	// Iterate over the map keys
	for k := range cache.indices {
		if cache.indices[k] > cache.indices[key] {
			cache.indices[k] -= 1
		}
	}
	// Delete key from cache.indices map
	delete(cache.indices, key)
	// Return the removed value
	return cache.data[cache.indices[key]]
}

// The Show function is used for getting the cache data.
// The function read locks the mutex then returns the
// cache.data. Once the function has returned, the
// mutex unlocks

// Show the cache
/* Returns 							*/
/* 	cache.data: []interface{} 		*/
func (cache *Cache) Show() []interface{} {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the cache values
	return cache.data
}

// The ShowIndexMap function is used for getting the cache
// slice indices. The function read locks the mutex
// then returns the cache.indices. Once the function has returned,
// the mutex unlocks

// Show the cache
/* Returns 								*/
/* 	cache.indices: map[string]int 		*/
func (cache *Cache) ShowIndexMap() map[string]int {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the cache indices map
	return cache.indices
}

// The ShowKeys function is used to get a slice of all
// the cache keys. The function read locks the cache mutex
// before iterating over the cache indices map, adding each
// of the keys to the keys slice.
//
// The function then returns the slice of keys. Once the
// function returns, the cache mutex is unlocked
//
// Returns all the cache keys in a slice
//
/* Returns 							*/
/* 	keys: []string 					*/
func (cache *Cache) ShowKeys() []string {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Define Variables
	var (
		// keys -> The slice containing the keys
		keys []string = make([]string, len(cache.indices))
		// i -> Track the index for setting the keys
		i int = 0
	)
	// Iterate over the cache indices map
	for k := range cache.indices {
		keys[i] = k
		i++
	}
	return keys
}

// The Flush function is used to clear the cache
// data and the cache indices. The function locks
// the cache mutex before resetting the cache data
// and the cache indices. Once the function returns
// the cache mutex is unlocked

// Clear the cache data
func (cache *Cache) Flush() {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Reset the cache variables
	cache.data = []interface{}{}
	cache.indices = map[string]int{}
}
