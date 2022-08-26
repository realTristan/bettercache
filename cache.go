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
	} else
	// Unlimited size
	{
		return &Cache{
			Data:  []byte{'*'},
			Mutex: &sync.RWMutex{},
		}
	}
}

// Initialize Cache
func Init(size int) *Cache {
	// Set global variables
	MaxCacheSize = size

	// Create new cache
	var c *Cache = _Init_(size)
	return c
}

// The serialize() function covnerts the byte cache into
// a map that can be used for reading keys, deleting keys, etc.
func (cache *Cache) serialize() map[string]map[string]string {
	// Convert the byte cache into a json
	// serializable string
	var c = []byte{'{'}
	c = append(c, cache.Data[1:len(cache.Data)-1]...)
	c = append(c, '}')

	// Unmarshal the serialized cache
	var tmp map[string]map[string]string
	json.Unmarshal(c, &tmp)

	// Return the map
	return tmp
}

// The Set() function sets the value for the
// provided key inside the cache.
//
// Example: key1: map[string]string{"1": "2"}
func (cache *Cache) Set(key string, data map[string]string) {
	// Lock/Unlock the mutex
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Marhsal the data
	var tmp, _ = json.Marshal(map[string]map[string]string{
		key: data,
	})
	// Set the byte cache value
	cache.Data = append(
		cache.Data, append(tmp[1:len(tmp)-1], ',')...)
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

// The ExpireKey() function removes the provided key
// from the cache after the given time
func (cache *Cache) ExpireKey(key string, _time time.Duration) {
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

// The SerializeData() function takes all the cache data
// and adds it's key's values to it's own map
//
// Example: {"key1": {"a": "b", "1": "2"}}
// Will be converted to: {"a": "b", "1", "2"}
func (cache *Cache) SerializeData() map[string]string {
	// Define variables
	var (
		// res -> Result map
		res map[string]string = make(map[string]string)
		// _cache -> Serialized cache
		_cache = cache.serialize()
	)
	// Iterate over the serialized cache
	for _, v := range _cache {
		for k, _v := range v {
			res[k] = _v
		}
	}
	// Return the result
	return res
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
func (cache *Cache) DumpData() map[string]map[string]string {
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

// The Remove() function locks then unlocks the
// cache data to ensure safety before serializing the
// byte cache into a map
//
// Once the cache is converted into a map, it deletes
// the key from said map then re-converts the map into
// a byte slice, setting the cache.Data to said slice
func (cache *Cache) Remove(key string) {
	// Lock/Unlock the mutex
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Get the cache map and delete the key
	var _cache = cache.serialize()
	delete(_cache, key)

	// Set the Cache Data
	var data, _ = json.Marshal(_cache)
	cache.Data = append([]byte{'*'}, data[1:]...)
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
	for i := 0; i < len(TempCache); i++ {
		// Break the loop if over the text search limit
		if TS.Limit > 0 && len(Result) >= TS.Limit {
			break
		} else

		// Check whether the current index is
		// in a string or not
		if i > 0 {
			if TempCache[i] == '"' && TempCache[i-1] != '\\' {
				inString = !inString
			}
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
	for i := 0; i < len(cache.Data); i++ {
		// Check whether the current index is
		// in a string or not
		if i > 0 {
			if cache.Data[i] == '"' && cache.Data[i-1] != '\\' {
				inString = !inString
			}
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
