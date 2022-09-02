# Better Cache ![Stars](https://img.shields.io/github/stars/realTristan/bettercache?color=brightgreen) ![Watchers](https://img.shields.io/github/watchers/realTristan/bettercache?label=Watchers)
![banner (5)](https://user-images.githubusercontent.com/75189508/186757681-6b7f97e8-ec37-448a-83cc-75106ed16309.png)

Modern Caching System for Go.

# About
- Better Cache is a Modern, Native golang caching system that utilizes slices for storing data. Because of this, we can perform lightning fast full text searches, full text removes, etc.

# Benchmarks

<h3>Full Text Search</h3>

```
    (1) Cache Size: 25 -> ~779ns
    (10) Cache Size: 250 -> ~3.203µs
    (100) Cache Size: 2,500 -> ~13.324µs
    (1,000) Cache Size: 25,000 -> ~160.049µs
    (10,000) Cache Size: 250,000 -> ~2.059ms
```

# Quick Usage

```go
package main

// Import Packages
import (
    "fmt"
    cache "github.com/realTristan/bettercache"
)

func main() {
    // Initialize the cache
    var Cache *cache.Cache = cache.Init(-1) // -1 (no size limit)

    // Add key1 to the cache
    Cache.Set(&cache.SetData{
        Key:      "key1",       // The cache key
        Value:    "value1",     // The cache value
        FullText: true,         // If true, Value converts to a string
    })

    // Get key from the cache
    var data string = Cache.Get("key1")
    fmt.Println(data)

    // Full Text Search for the key's contents
    var res []string = Cache.FullTextSearch(&cache.TextSearch{
        Limit:      -1,                 // No limit
        Query:      []byte("value"),    // Search for "value"
        StrictMode: false,              // Ignore CAPS
    })
    fmt.Println(res)

    // Remove key1 from the cache
    var removedKey string = Cache.Remove("key1")
    fmt.Println(removedKey)
}
```

# Functions

```go

// The _Cache struct has six primary keys
// CurrentSize: int { "The current map size" } 
// maxSize: int { "The maximum map size" } 
// mutex: *sync.RWMutex { "The mutex for locking/unlocking the data" } 				  
// mapData: map[interface{}]interface{} { "The Main Data Cache Values" } 								  
// fullTextData: []string { "The Full Text Data Cache Values" } 					  
// fulltextIndices: map[string]int { "The Cache Keys holding the full text indices of the Cache Values" } 	
type _Cache struct {
	currentSize     int
	maxSize         int
	mutex           *sync.RWMutex
	fullTextData    []string
	fullTextIndices map[interface{}]int
	mapData         map[interface{}]interface{}
}

// The SetData struct has three primary keys
// Key: string { "The Cache Key" }		                            
// Value: string { "The Cache Value" }	                            
// FullText: string { "Whether to enable full text functions " }	
// WARNING
// If FullText is set to true, it converts the Value to a string	
type SetData struct {
	Key      interface{}
	Value    interface{}
	FullText bool
}

// The TextRemove struct has two primary keys
// Query: string { "The string to find within values" } 				
// Amount: int { "The amount of keys to remove (Set to -1 for all)" }	
type TextRemove struct {
	Query  string
	Amount int
}

// The TextSearch struct contains five primary keys
// Query: string { "The string to search for in the cache values" } 		
// Limit: int { "The amount of search results. (Set to -1 for no limit)" }  
// StrictMode: bool { "Set to false to ignore caps in query comparisons" }  
// PreviousQueries: map[string][]string { "The Previous Searches" } 			
type TextSearch struct {
	Query               string
	Limit               int
	StrictMode          bool
	PreviousQueries     map[string][]string
}

// The Full Text Search function is used to find all cache values
// that contain the provided query.
//
// The Full Text Search function iterates over the cache slice
// and uses the strings.Contains function to check whether
// the cache value contains the query. If the value contains the
// query, it will append the cache value to the { res: []string }
// slice. Once the cache has been fully iterated over, the function
// will return the { res: []string } slice.
//
// If the user is not using strictmode it will set the cache value
// and the provided query to lowercase

// Performs a full text search using the cache values
// Parameters 
// 	TS: *TextSearch = &TextSearch{
//		Query               	string
//		Limit               	int
//		StrictMode          	bool
//		StorePreviousSearch 	bool
//		PreviousSearch      	map[string][]string
//})

//
// If you want to store the previous text search you made, you can set the
// StorePreviousSearch to true. This will set the key in the previous search
// to the provided TextSearch.Query and the value to the result slice.
//
// >> Returns 			
// res: []string	 	
func (cache *Cache) FullTextSearch(TS *TextSearch) []string {}

// The Full Text Remove function is used to find all cache values
// that contain the provided query and remove their keys from the cache
//
// The Full Text Remove function iterates over the cache slice
// and uses the strings.Contains function to check whether
// the cache value contains the query. If the value contains the
// query, it will append the cache value to the { res: []string }
// slice, then remove the key from the cache.
// Once the cache has been fully iterated over, the function
// will return the { res: []string } slice.

// Removes keys in the cache depending on whether their values
// contain the provided query
// Parameters 
// 	TS: *TextRemove = &TextRemove{
//		Query               	string
//		Amount               	int
//})
//
// If you want to remnove all the values, either call the FullTextRemoveAll()
// function or set the TextRemove.Amount to -1
//
// >> Returns       
// res: []string    
func (cache *Cache) FullTextRemove(TR *TextRemove) []string {}

// The Full Text Remove All function utilizes the Full Text Remove function
// to remove the keys whos values contain the provided query.

// Removes all cache keys that contain the provided query in their values
// Paramters 
//	query: string { "The string to query for" } 
func (cache *Cache) FullTextRemoveAll(query string) []string {}

// The Init function is used for initalizing the Cache struct
// and returning the newly created cache object.
// You can create the cache by yourself, but using the Init()
// function is much easier

// Initializes the cache object
// Parameters: 												
// 	size: int { "The Size of the cache map and slice" }  	
//
// Returns 													
// 	cache: *Cache 											
func Init(size int) *Cache {}

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
// Parameters: 
// 	SD &SetData = *SetData{
//		Key: interface{},
//		Value: interface{},
//		FullText: bool,
//} 
func (cache *Cache) Set(SD *SetData) {}

// The ExistsinFullText function is used for checking whether a key
// exists in the full text cache or not. The function read locks
// the cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
// Parameters: 								
// 	key: interface{} { "The Cache Key" } 	
//
// Returns 									
// 	doesExist: bool 						
func (cache *Cache) ExistsInFullText(key interface{}) bool {}

// The ExistsInMap function is used for checking whether a key
// exists in the main cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
// Parameters: 							
// 	key: interface{} { "The Cache Key" } 	
//
// Returns 								
// 	doesExist: bool 					
func (cache *Cache) ExistsInMap(key interface{}) bool {}

// The Exists function is used for checking whether a key
// exists in the cache or not. The function read locks the
// cache mutex before returning whether the key is in the
// cache. Once the function returns, the mutex is unlocked

// Returns whether the provided key exists in the cache
// Parameters: 								
// 	key: interface{} { "The Cache Key" } 	
//
// Returns 									
// 	doesExist: bool 						
func (cache *Cache) Exists(key interface{}) bool {}

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
// Parameters: 								
// 	key: interface{} { "The Cache Key" } 	
//
// Returns 									
// 	cacheValue: interface{} 				
func (cache *Cache) Get(key interface{}) interface{} {}

// The Remove function is used to remove a value from the cache
// using it's corresponding key. The function full locks the mutex
// before modifying the cache.mainData, removing the cache value.
//
// Once the cache.mapData value is removed, the function moves onto
// iterating through the cache.mainIndices map reducing all cache indices
// that are post-the-removed-value. The function then returns the
// removed value. Once the function returns, the cache mutex is unlocked.

// Removes a key from the cache
// Parameters: 								
// 	key: interface{} { "The Cache Key" } 	
//
// Returns 									
// 	removedValue: interface{} 				
func (cache *Cache) Remove(key interface{}) interface{} {}

// The Show function is used for getting the cache data.
// The function read locks the mutex then returns the
// cache.mapData and the cache.fullTextData,
// Once the function has returned, the mutex unlocks

// Show the cache
// Returns 								
// 	cache.mainData: []interface{} 		
func (cache *Cache) Show() (map[interface{}]interface{}, []string) {}

// The ShowFTIndices function is used for getting the cache
// slice indices. The function read locks the mutex
// then returns the cache.fullTextIndices. Once the function has returned,
// the mutex unlocks

// Show the cache
// Returns 								
// 	cache.mainIndices: map[interface{}]int 		
func (cache *Cache) ShowFTIndices() map[interface{}]int {}

// The ShowKeys function is used to get a slice of all
// the cache keys. The function read locks the cache mutex
// before iterating over the cache indices map, adding each
// of the keys to the keys slice.
//
// The function then returns the slice of keys. Once the
// function returns, the cache mutex is unlocked

// Returns all the cache keys in a slice
//
// Returns 							
// 	keys: []interface{}				
func (cache *Cache) ShowKeys() []interface{} {}

// The Clear function is used to clear the cache
// data and the cache indices. The function locks
// the cache mutex before resetting the cache data
// and the cache indices. Once the function returns
// the cache mutex is unlocked

// Clear the cache data
func (cache *Cache) Clear() *Cache {}

// The GetMaxSize function is used to get the maximum size
// of the cache. The function read locks the cache mutex
// before returning the cache maxSize. Once the function
// returns, the cache mutex is unlocked.

// Returns the caches maximum size (int)
func (cache *Cache) GetMaxSize() int {}

// The GetCurrentSize function is used to get the current size
// of the cache. The function read locks the cache mutex
// before returning the cache currentSize. Once the function
// returns, the cache mutex is unlocked.

// Return the cache current size (int)
func (cache *Cache) GetCurrentSize() int {}


// The Flush function locks the mutex before calling
// the flushToFile function. Once the function has been
// called and the Flush function returns, the mutex is unlocked

// The Flush function is used to write the cache
// data to a BetterCache file
//
// Paramters: 
// 	path: string { "The path to the BetterCache file" } 
func (cache *Cache) Flush(path string) {}

// The GetPreviousQueries function is used to return the
// slice of values for a previous query
//
// Paramters 
// query: string { "The Previous Query" } 
//
// Returns 
// results: []string { "The Query Results" } 
func (TS *TextSearch) GetPreviousQueries(query string) []string {}

// The GetPreviousQueries function is used to delete a
// previous query from the PreviousQueries map
//
// Paramters 
// query: string { "The Previous Query" } 
func (TS *TextSearch) DeletePreviousQuery(query string) {}

// The ClearPreviousQueries function is used to reset
// the previous queries map
//
// Paramters 
// size: int { "The Size of PreviousQueries Map (Set to -1 for no limit)" } 
func (TS *TextSearch) ClearPreviousQueries(size int) {}

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
