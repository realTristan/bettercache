package cache

// The TextSearch struct contains three primary keys
/* Query: []byte -> What to query for									*/
/* StrictMode: bool -> Whether to convert the cache data to lowercase	*/
/* Limit: int -> The number of results to return						*/
type TextSearch struct {
	Query      []byte
	StrictMode bool
	Limit      int
}

// The isLetter() function returns whether the provided
// character is a letter or not.
func isLetter(c byte) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

// The equalCharAtIndex() returns whether the characters
// at the given index are equal to eachother
//
// If not StrictMode then it doesn't matter whether the
// character is uppercase/lowercase
func (cache *Cache) equalCharAtIndex(i int, index int, TS *TextSearch) bool {
	// Ensure the character is a letter and check
	// Whether strict mode is disabled
	if !TS.StrictMode && isLetter(cache.Data[i]) {
		return cache.Data[i]-32 == TS.Query[index] ||
			cache.Data[i]+32 == TS.Query[index] ||
			cache.Data[i] == TS.Query[index]
	}
	// Return the strict mode result
	return cache.Data[i] == TS.Query[index]
}

// The FullTextSearch() function iterates through the
// cache data and returns the json value of a key.
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
		// valueIndex -> Track the index of the TS.Query
		// This is used for checking whether the Query is in
		// the middle of the cache key start and the cache key end
		valueIndex int = -1
		// index -> Track the index of the TS.Query
		index int = 0
	)

	// Iterate over the lowercase cache string
	for i := 1; i < len(cache.Data); i++ {

		// Break the loop if over the text search limit
		if TS.Limit > 0 && len(Result) >= TS.Limit {
			break
		}

		// Check if the strings are equal
		if cache.equalCharAtIndex(i, index, &TS) {
			index++
		} else {
			// Reset the index
			index = 0
		}
		// Check if the key is present and the current
		// index is standing at that key
		if index == len(TS.Query) && valueIndex < 0 {
			valueIndex = i
			index = 0
		}

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
		// Also Make sure the map start has been established
		if cache.Data[i] == '}' && mapStart > 0 {

			// Check whether the query index is in between the map start
			// and the map end
			if valueIndex > mapStart && valueIndex < i+1 {

				// Append the data to the result array
				Result = append(Result, string(cache.Data[mapStart+1:i]))
			}

			// Reste the indexing variables
			mapStart, valueIndex = -1, -1
			dataLength = []byte{}
		}
	}
	// Return the result
	return Result
}
