package main

// Import modules
import (
	"bytes"
	"encoding/json"
	"fmt"
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

// The serialize() function covnerts the byte cache into
// a map that can be used for reading keys, deleting keys, etc.
func (cache *Cache) serialize() map[string]string {
	// Convert the byte cache into a json
	// serializable string
	var c = []byte{'{'}
	c = append(c, cache.Data[1:len(cache.Data)-1]...)
	c = append(c, '}')

	// Unmarshal the serialized cache
	var tmp map[string]string
	json.Unmarshal(c, &tmp)

	// Return the map
	return tmp
}

// The Exists() function returns whether the
// provided key exists in the cache
func (cache *Cache) Exists(key string) bool {
	var _, i = cache.serialize()[key]
	return i
}

// The GetByteSize() function returns the current size of the
// cache bytes and the cache maximum size
func (cache *Cache) GetByteSize() (int, int) {
	return len(cache.Data), MaxCacheSize
}

// The GetMapSize() function returns the
// amount of keys in the cache map and the cache
// maximum size
func (cache *Cache) GetMapSize() (int, int) {
	return len(cache.serialize()), MaxCacheSize
}

// The Expire() function removes the provided key
// from the cache after the given time
func (cache *Cache) Expire(key string, _time time.Duration) {
	go func(key string, _time time.Duration) {
		time.Sleep(_time)
		cache.Remove(key)
	}(key, _time)
}

// The Flush() function resets the cache
// data. Make sure to use this function when
// clearing the cache!
func (cache *Cache) Flush() {
	cache.Data = []byte{'*'}
}

// The DumpBytes() function returns the cache
// bytes. Use the DumpData() function for returning
// the actual map
func (cache *Cache) DumpBytes() []byte {
	return cache.Data
}

// The DumpJson() function returns the cache
// as a json map.
func (cache *Cache) DumpJson() string {
	return string(
		append([]byte{'{'},
			append(cache.Data[1:len(cache.Data)-1], '}')...))
}

// The DumpData() function returns the serialized
// cache map. Use the DumpData() function for returning
// the cache bytes
func (cache *Cache) DumpData() map[string]string {
	return cache.serialize()
}

// The GetKeys() function returns all the keys
// inside the cache
func (cache *Cache) GetKeys() []string {
	var res []string
	for k := range cache.serialize() {
		res = append(res, k)
	}
	return res
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

	// Set the new key
	key = fmt.Sprintf(`"%s":{`, key)

	// Lock/Unlock the mutex
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Set the byte cache value
	cache.Data = append(
		cache.Data, append([]byte(key),
			append([]byte(data), []byte{'}', ','}...)...)...)

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
		// inString -> Track whether bracket is inside a string
		inString bool = false
		// mapStart -> Track opening bracket
		mapStart int = -1
		// closeBracketCount -> Track closing brackets per map
		closeBracketCount int = 0
		// Result -> Array with all maps containing the query
		Result []string
		// Set the temp cache
		TempCache []byte = cache.Data
	)

	// Check if strict mode is enabled
	// If true, convert the temp cache to lowercase
	if !TS.StrictMode {
		TempCache = bytes.ToLower(cache.Data)
	}

	// Iterate over the lowercase cache string
	for i := 1; i < len(TempCache); i++ {
		// Break the loop if over the text search limit
		if TS.Limit > 0 && len(Result) >= TS.Limit {
			break
		} else

		// Check whether the current index is
		// in a string or not
		if TempCache[i] == '"' && TempCache[i-1] != '\\' {
			inString = !inString
		}

		// Check if current index is the start of a map
		if TempCache[i] == '{' && !inString {
			if mapStart == -1 {
				mapStart = i
			}
			closeBracketCount++
		} else

		// Check if the current index is the end of the map
		if TempCache[i] == '}' && !inString {
			if closeBracketCount == 1 {
				// Check if the map contains the query string
				if bytes.Contains(TempCache[mapStart:i+1], TS.Query) {
					// Append the json map to the result array
					Result = append(Result, string(cache.Data[mapStart:i+1]))
				}
				// Reset indexing variables
				closeBracketCount = 0
				mapStart = -1
			} else {
				closeBracketCount--
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
	// Set the new key
	key = fmt.Sprintf(`"%s":{`, key)

	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Define Variables
	var (
		// inString -> Track whether bracket is inside a string
		inString bool = false
		// startIndex -> Track the start of the key value
		startIndex int = -1
		// index -> Track the key indexes
		index int = 0
	)
	// Iterate over the lowercase cache string
	for i := 1; i < len(cache.Data); i++ {
		// Check whether the current index is
		// in a string or not
		if cache.Data[i] == '"' && cache.Data[i-1] != '\\' {
			inString = !inString
		}

		// Check if current index is the start of a map
		if index == len(key) {
			startIndex = i - 1
			index = 0
		} else if cache.Data[i] == key[index] {
			if startIndex < 0 {
				index++
			}
		} else {
			index = 0
		}
		// Check if the current index is the end of the map
		if cache.Data[i] == '}' && !inString {
			if startIndex > 0 {
				return string(append(cache.Data[startIndex:i], '}'))
			}
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
	// Set the new key
	key = fmt.Sprintf(`"%s":{`, key)

	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Define Variables
	var (
		// inString -> Track whether bracket is inside a string
		inString bool = false
		// startIndex -> Track the start of the key value
		startIndex int = -1
		// index -> Track the key indexes
		index int = 0
	)
	// Iterate over the lowercase cache string
	for i := 1; i < len(cache.Data); i++ {
		// Check whether the current index is
		// in a string or not
		if cache.Data[i] == '"' && cache.Data[i-1] != '\\' {
			inString = !inString
		}

		// Check if current index is the start of a map
		if index == len(key) {
			startIndex = i - len(key)
			index = 0
		} else if cache.Data[i] == key[index] {
			if startIndex < 0 {
				index++
			}
		} else {
			index = 0
		}
		// Check if the current index is the end of the map
		if cache.Data[i] == '}' && !inString {
			if startIndex > 0 {
				// Store the removed value
				var data string = string(append(cache.Data[startIndex:i], '}'))
				// Remove the value
				cache.Data = append(cache.Data[:startIndex], cache.Data[i+2:]...)
				// Return the value removed
				return data[len(key)-1:]
			}
		}
	}
	// Return empty string
	return ""
}
