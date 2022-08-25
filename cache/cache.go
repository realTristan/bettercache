package cache

// Import modules
import (
	"bytes"
	"encoding/json"
	"sync"
)

// TextSearch object
type TextSearch struct {
	Query      []byte
	StrictMode bool
	Limit      int
}

// Cache struct
type Cache struct {
	Data  []byte
	Mutex *sync.RWMutex
}

// Initialize Cache
func Init() *Cache {
	return &Cache{
		Data:  []byte{'*'},
		Mutex: &sync.RWMutex{},
	}
}

// The Set() function sets the data for the
// given key in the cache
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

// Convert the cache into a map
func (cache *Cache) Serialize() map[string]map[string]string {
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

// Remove a key from the cache
func (cache *Cache) Remove(key string) {
	// Lock/Unlock the mutex
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Get the cache map and delete the key
	var _cache = cache.Serialize()
	delete(_cache, key)

	// Set the Cache Data
	var data, _ = json.Marshal(_cache)
	cache.Data = append([]byte{'*'}, data[1:]...)
}

// Get a key from the cache
func (cache *Cache) Get(key string) map[string]string {
	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Return the cache data
	return cache.Serialize()[key]
}

// The FullTextSearch() function..
func (cache *Cache) FullTextSearch(TS TextSearch) []map[string]interface{} {
	// Lock/Unlock the mutex
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()

	// Define Variables
	var (
		// Track whether bracket is inside a string
		inString bool = false
		// courseMapStart -> Track opening bracket
		courseMapStart int = -1
		// closeBracketCount -> Track closing brackets per course map
		closeBracketCount int = 0
		// similarResult -> Array with all courses that contain the query
		Result []map[string]interface{}
		// Set the temp cache
		TempCache []byte = cache.Data
	)

	// Check if strict mode is enabled
	if !TS.StrictMode {
		TempCache = bytes.ToLower(cache.Data)
	}

	// Iterate over the lowercase cache string
	for i := 0; i < len(TempCache); i++ {
		// Break the loop if there's too many similar courses
		if TS.Limit > 0 && len(Result) > TS.Limit {
			break
		}

		// Check whether in string or not
		if i > 0 {
			if TempCache[i] == '"' && TempCache[i-1] != '\\' {
				inString = !inString
			}
		}

		// Check if current index is the start of
		// the course data map
		if TempCache[i] == '{' && !inString {
			if courseMapStart == -1 {
				courseMapStart = i
			}
			closeBracketCount++
		} else

		// Check if the current index is the end of
		// the course data map
		if TempCache[i] == '}' && !inString {
			if closeBracketCount == 1 {
				// Check if the map contains the query string
				if bytes.Contains(TempCache[courseMapStart:i+1], TS.Query) {
					// Convert the string to a map
					var data map[string]interface{}
					json.Unmarshal(cache.Data[courseMapStart:i+1], &data)

					// Append the map to the result array
					Result = append(Result, data)
				}
				// Reset indexing variables
				closeBracketCount = 0
				courseMapStart = -1
			} else {
				closeBracketCount--
			}
		}
	}
	// Return the combined arrays
	return Result
}
