package main

import (
	"bytes"
	"strings"
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
