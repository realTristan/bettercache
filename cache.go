package main

// Import modules
import (
	"bytes"
	"encoding/json"
	"sync"
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
/* Data: []byte -> The Cache Data 							 */
/* Mutex: *sync.Mutex -> Used for locking/unlocking the data */
type Cache struct {
	Data  []byte
	Mutex *sync.RWMutex
}

// Initialize Cache
func Init(size int) *Cache {
	var c *Cache = &Cache{
		Data:  make([]byte, size+1),
		Mutex: &sync.RWMutex{},
	}
	c.Data[0] = '*'
	return c
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
	var (
		tmp, _ = json.Marshal(map[string]map[string]string{
			key: data,
		})
	)
	// Set the byte cache value
	cache.Data = append(
		cache.Data, append(tmp[1:len(tmp)-1], ',')...)
}

// The serialize() function covnerts the byte cache into
// a map that can be used for reading keys, deleting keys, etc.
func (cache *Cache) serialize() map[string]map[string]string {
	// Convert the byte cache into a json
	// serializable string
	var _cache = []byte{'{'}
	_cache = append(_cache, cache.Data[1:len(cache.Data)-1]...)
	_cache = append(_cache, '}')

	// Unmarshal the serialized cache
	var tmp map[string]map[string]string
	json.Unmarshal(_cache, &tmp)

	// Return the map
	return tmp
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

// The Get() function read locks then read unlocks
// the cache data to ensure safety before serializing
// the byte cache to a map.
//
// Once the cache is converted into a map, it will then
// return the value of the provided key
func (cache *Cache) Get(key string) map[string]string {
	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Return the cache data
	return cache.serialize()[key]
}

// The FullTextSearch() function iterates through the cache data
// and returns the value of a key. This value contains the Query
// defined in the provided TextSearch object
//
// To ensure safety, the cache data is locked then unlocked once
// no longer being used
func (cache *Cache) FullTextSearch(TS TextSearch) []map[string]string {
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
		Result []map[string]string
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
					// Convert the string to a map
					var data map[string]string
					json.Unmarshal(cache.Data[mapStart:i+1], &data)

					// Append the map to the result array
					Result = append(Result, data)
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
