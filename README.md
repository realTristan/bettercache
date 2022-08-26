# Better Cache ![Stars](https://img.shields.io/github/stars/realTristan/BetterCache?color=brightgreen) ![Watchers](https://img.shields.io/github/watchers/realTristan/BetterCache?label=Watchers)
![banner (5)](https://user-images.githubusercontent.com/75189508/186757681-6b7f97e8-ec37-448a-83cc-75106ed16309.png)

Lightning Fast Caching System for Go.

# About
- Better Cache is an ultra fast caching system that uses an array of bytes for storing data instead of the very common map caching system. Because the data is stored in an array of bytes instead of a map, it enables lightning fast full text searches. (within microseconds to milliseconds)
- Better Cache also uses solely native golang modules which makes it fast, lightweight and secure.

# Benchmarks

<h3>Full Text Search</h3>

```
    (1) Cache Size: 25 -> ~1.437µs
    (10) Cache Size: 250 -> ~12.246µsµs
    (100) Cache Size: 2,500 -> ~42.724µs
    (1,000) Cache Size: 25,000 -> ~267.382µs
    (10,000) Cache Size: 250,000 -> ~1.934ms
```

# Quick Usage

```go
package main

import (
    "fmt"
    cache "github.com/realTristan/BetterCache"
)

func main() {
    // Initialize the cache
	var Cache *cache.Cache = cache.Init(100) // 100 bytes

    // Add key1 to the cache
    Cache.Set("key1", map[string]string{
		"summary": "My name is \"Tristan\"",
	})

    // Get key from the cache
    var data = Cache.Get("key1")
    fmt.Println(data)

    // Full Text Search for the key's contents
	var res = Cache.FullTextSearch(Cache.TextSearch{
		Limit:      -1,                 // No limit
		Query:      []byte("tristan"),  // Search for "tristan"
		StrictMode: false,              // Ignore CAPS
	})
    fmt.Println(res)

    // Remove key1 from the cache
    Cache.Remove("key1")
}

```

# Functions

```go
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
```

# License
MIT License

Copyright (c) 2022 Tristan Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
