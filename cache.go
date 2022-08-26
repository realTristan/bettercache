package main

// Import modules
import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"
)

// The TextSearch struct contains three primary keys
/* Query: []byte -> What to query for									*/
/* StrictMode: bool -> Whether to convert the cache data to lowercase	*/
/* Limit: int -> The number of results to return						*/
type TextSearch struct {
	Query      []byte
	StrictMode bool
	Limit      int
}

// The Cache struct contains two primary keys
/* Data: []byte -> The Cache Data in Bytes						 	 */
/* Mutex: *sync.Mutex -> Used for locking/unlocking the data 	 	 */
type Cache struct {
	Data  []byte
	Mutex *sync.RWMutex
}

// Global Variables
var (
	// CacheSize -> The size of the cache
	MaxCacheSize int
)

// The _Init_() function returns a Cache object
// based off the provided cache size
func _Init_(size int) *Cache {
	// Limited Size
	if size > 0 {
		var c *Cache = &Cache{
			Data:  make([]byte, size+1),
			Mutex: &sync.RWMutex{},
		}
		c.Data[0] = '*'
		return c
	}
	// Unlimited size
	return &Cache{
		Data:  []byte{'*'},
		Mutex: &sync.RWMutex{},
	}

}

// The Init() function creates the Cache
// object depending on what was entered for
// the size of the cache
func Init(size int) *Cache {
	// Set global variables
	MaxCacheSize = size

	// Create new cache
	var c *Cache = _Init_(size)
	return c
}

// The Exists() function returns whether the
// provided key exists in the cache
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) Exists(key string) bool {
	return len(cache.Get(key)) > 0
}

// The GetByteSize() function returns the current size of the
// cache bytes and the cache maximum size
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) GetByteSize() (int, int) {
	return len(cache.Data), MaxCacheSize
}

// The Expire() function removes the provided key
// from the cache after the given time
//
// Suggested to run this function in your
// own monitored goroutine
func (cache *Cache) Expire(key string, _time time.Duration) {
	time.Sleep(_time)
	cache.Remove(key)
}

// The Flush() function resets the cache
// data. Make sure to use this function when
// clearing the cache!
func (cache *Cache) Flush() {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Reset the cache data
	cache.Data = []byte{'*'}
}

// The ShowBytes() function returns the cache bytes
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) ShowBytes() []byte {
	return cache.Data
}

// The Show() function returns the cache as a string
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) Show() string {
	return string(cache.Data)
}

// The Set() function sets the value for the
// provided key inside the cache.
//
// Example: {"key1": "my name is tristan!"},
//
// Returns the removed value of the previously
// defined key
func (cache *Cache) Set(key string, data string) string {
	var removedValue string = cache.Remove(key)

	// Lock/Unlock the mutex
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Set the byte cache value
	cache.Data = append(
		cache.Data, append([]byte(
			fmt.Sprintf(`|%s|:~%d{`, key, len(data))), []byte(fmt.Sprintf(`%s}`, data))...)...)

	// Return the removed value
	return removedValue
}

// The FullTextSearch() function iterates through the cache data
// and returns the json value of a key.
// This value contains the Query defined in the provided
// TextSearch object
//
// To ensure safety, the cache data is locked then unlocked once
// no longer being used
func (cache *Cache) FullTextSearch(TS TextSearch) []string {
	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Define Variables
	var (
		// Result -> Array with all maps containing the query
		Result []string
		// mapStart -> Track opening bracket
		mapStart int = -1
		// Track the length of the data
		dataLength = []byte{}

		// Check if strict mode is enabled
		// If true, convert the temp cache to lowercase
		isStrictMode = func() []byte {
			// Return the cache in lowercase
			if !TS.StrictMode {
				return bytes.ToLower(cache.Data)
			}
			// Return the cache as is
			return cache.Data
		}
	)

	// Iterate over the lowercase cache string
	for i := 1; i < len(cache.Data); i++ {

		// Break the loop if over the text search limit
		if TS.Limit > 0 && len(Result) >= TS.Limit {
			break
		} else

		// Check if current index is the start of a map
		if cache.Data[i] == '{' {
			// Make sure the map start has NOT been established
			if mapStart == -1 {
				mapStart = i
			}

			// Get the length of the cache data
			if len(dataLength) == 0 {
				dataLength = cache.Data[mapStart-2 : i]
			}
		} else

		// Check if the current index is the end of the map
		if cache.Data[i] == '}' {

			// Make sure the map start has been established
			if mapStart > 0 {

				// Check if the current data contains the query
				if bytes.Contains(isStrictMode()[mapStart:i+1], TS.Query) {

					// Append the data to the result array
					Result = append(Result, strings.ReplaceAll(
						string(cache.Data[mapStart+1:i]), "~|", ""))
				}

				// Reset the data length variable
				// and the map start variable
				dataLength = []byte{}
				mapStart = -1
			}
		}
	}
	// Return the result
	return Result
}

// The Get() function read locks then read unlocks
// the cache data to ensure safety before returning
// a json map with the key's value
func (cache *Cache) Get(key string) string {
	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Set the new key
	var newKey string = fmt.Sprintf(`|%s|:~`, key)

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
		if index == len(newKey) {

			// Set the starting index and reset the
			// key's index variable
			startIndex = i + 1
			index = 0
		} else

		// Check if the current index value is
		// equal to the newKey's current index
		if cache.Data[i] == newKey[index] {
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
		}
	}
	// Return empty string
	return ""
}

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

	// Set the new key
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
