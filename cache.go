package cache

// Import Packages
import (
	"fmt"
	"sync"
)

// The Cache struct has three primary keys (all unexported)
/* mutex: *sync.RWMutex { "The mutex for locking/unlocking the data" } 				  */
/* data: []string { "The Cache Values" } 											  */
/* ftData: []int { "The Indexes to use for full text functions" }					  */
/* index: map[string]int { "The Cache Keys holding the indexes of the Cache Values" } */
type Cache struct {
	mutex *sync.RWMutex
	data  []interface{}
	index map[string]int
}

// The SetData struct has three primary keys
/* */
type SetData struct {
	Key      string
	Value    interface{}
	FullText bool
}

// The Init function is used for initalizing the Cache struct
// and returning the newly created cache object.
// You can create the cache by yourself, but using the Init()
// function is much easier
//
/* >> Parameters 										*/
/* size: int { "The Size of the cache map and slice" }  */
//
/* >> Returns 			*/
/* cache: *Cache 		*/
func Init(size int) *Cache {
	// If the user passed a size variable greater
	// tham zero.
	if size > 0 {
		// If so, it will use the make() function
		// for created a data slice and index map
		// with a set size
		return &Cache{
			mutex: &sync.RWMutex{},
			data:  make([]interface{}, size),
			index: make(map[string]int, size),
		}
	}
	// Return a cache with no limit to it's
	// size. It is recommended that you provide a size.
	return &Cache{
		mutex: &sync.RWMutex{},
		data:  []interface{}{},
		index: make(map[string]int),
	}
}

// The Set function is used for setting a new value inside
// the cache data. The Set function locks the cache mutex to
// prevent data overwriting before checking whether the provided
// key already exists. If it does, it will call the Remove() function
// to remove that key from the cache.
//
// Once finished with the removal process, the function modifies
// the cache data, adding the value to the slice and adds the slice
// index along with the key to the cache index map
//
/* >> Parameters 										 			    */
/* key: string { "The Cache Key" }  					 			    */
/* value: interface{} { "The value to set the key to" }  			    */
/* fullText: bool { "Whether to add the value to the full text slice" } */
func (cache *Cache) Set(SD *SetData) {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// If key exists
	if _, t := cache.index[SD.Key]; t {
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

	// Set the key in the cache index map to the
	// index the key value is at.
	cache.index[SD.Key] = len(cache.data)

	// If the user set the AddToFullText to true
	// then add the cache data index to the ftData slice
	if SD.FullText {
		// Add the value into the cache data slice
		cache.data = append(cache.data,
			fmt.Sprintf("FT(true):%s:%v", SD.Key, SD.Value))
	} else {
		// Add the value into the cache data slice
		cache.data = append(cache.data, SD.Value)
	}
}

// Check if key exists
func (cache *Cache) Exists(key string) bool {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Check if the key exists within the cache.index
	if _, t := cache.index[key]; t {
		// It does
		return true
	}
	// It does not
	return false
}

// Get a cache key
func (cache *Cache) Get(key string) interface{} {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Make sure the key exists before returning
	// the key's value
	// If you don't check whether the key exists before
	// then it will return the key prior to the given in
	// the cache.index map
	if cache.Exists(key) {
		// Return the key value
		return cache.data[cache.index[key]]
	}
	// Return empty string
	return ""
}

// Set a cache key
func (cache *Cache) Remove(key string) {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Remove key from the cache slice
	cache.data = append(cache.data[:cache.index[key]],
		cache.data[cache.index[key]+1:]...)

	// Iterate over the map keys
	for k := range cache.index {
		if cache.index[k] > cache.index[key] {
			cache.index[k] -= 1
		}
	}
	// Delete key from cache.index map
	delete(cache.index, key)
}

// Return the cache values
func (cache *Cache) Show() []interface{} {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the cache values
	return cache.data
}

// Returns the index map
func (cache *Cache) ShowIndexMap() map[string]int {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Return the cache index map
	return cache.index
}

// Return the cache keys
func (cache *Cache) ShowKeys() []string {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Define Variables
	var (
		// keys -> The slice containing the keys
		keys []string = make([]string, len(cache.index))
		// i -> Track the index for setting the keys
		i int = 0
	)
	// Iterate over the cache index map
	for k := range cache.index {
		keys[i] = k
		i++
	}
	return keys
}

// Flush the cache data and the cache indexes
func (cache *Cache) Flush() {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Reset the cache variables
	cache.data = []interface{}{}
	cache.index = map[string]int{}
}
