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
// The Set() function sets the value for the
// provided key inside the cache.
//
// Example: {"key1": "my name is tristan!"},
//
// Returns the removed value of the previously
// defined key
func (cache *Cache) Set(key string, data string) string {}

// The Get() function read locks then read unlocks
// the cache data to ensure safety before returning
// a json map with the key's value
func (cache *Cache) Get(key string) string {}

// The Remove() function locks then unlocks the
// cache data to ensure safety before iterating through
// the cache bytes to look for the provided key
//
// once the key is found it'll search for it's closing
// bracket then remove the key from the cache bytes
//
// It will return the removed value
func (cache *Cache) Remove(key string) string {}

// The Exists() function returns whether the
// provided key exists in the cache
func (cache *Cache) Exists(key string) bool {}

// The GetByteSize() function returns the current size of the
// cache bytes and the cache maximum size
func (cache *Cache) GetByteSize() (int, int) {}

// The GetMapSize() function returns the
// amount of keys in the cache map and the cache
// maximum size
func (cache *Cache) GetMapSize() (int, int) {}

// The Expire() function removes the provided key
// from the cache after the given time
func (cache *Cache) Expire(key string, _time time.Duration) {}

// The Flush() function resets the cache
// data. Make sure to use this function when
// clearing the cache!
func (cache *Cache) Flush() {}

// The DumpBytes() function returns the cache
// bytes. Use the DumpData() function for returning
// the actual map
func (cache *Cache) DumpBytes() []byte {}

// The DumpJson() function returns the cache
// as a json map.
func (cache *Cache) DumpJson() string {}

// The DumpData() function returns the serialized
// cache map. Use the DumpData() function for returning
// the cache bytes
func (cache *Cache) DumpData() map[string]string {}

// The GetKeys() function returns all the keys
// inside the cache
func (cache *Cache) GetKeys() []string {}

// The FullTextSearch() function iterates through the cache data
// and returns the json value of a key.
// This value contains the Query defined in the provided
// TextSearch object
//
// To ensure safety, the cache data is locked then unlocked once
// no longer being used
func (cache *Cache) FullTextSearch(TS TextSearch) []string {}
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
