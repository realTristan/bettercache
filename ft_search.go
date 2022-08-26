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
//
// i: cache.data index
// n: query index
// q: query bytes
// sm: strict mode
func (cache *Cache) equalCharAtIndex(i int, n int, q []byte, sm bool) bool {
	// Ensure the character is a letter and check
	// Whether strict mode is disabled
	if !sm && isLetter(cache.data[i]) {
		return cache.data[i]-32 == q[n] ||
			cache.data[i]+32 == q[n] ||
			cache.data[i] == q[n]
	}
	// Return the strict mode result
	return cache.data[i] == q[n]
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
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// Define Variables
	var (
		// Result -> Array with all maps containing the query
		Result []string
		// mapStart -> Track opening bracket
		mapStart int = -1
		// valueIndex -> Track the index of the TS.Query
		// This is used for checking whether the Query is in
		// the middle of the cache key start and the cache key end
		valueIndex int = -1
		// index -> Track the index of the TS.Query
		index int = 0
	)

	// Iterate over the lowercase cache string
	for i := 1; i < len(cache.data); i++ {

		// Break the loop if over the text search limit
		if TS.Limit > 0 && len(Result) >= TS.Limit {
			break
		}

		// Check if the strings are equal
		if cache.equalCharAtIndex(i, index, TS.Query, TS.StrictMode) {
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
		if cache.data[i] == '{' {
			// Make sure the map start has NOT been established
			if mapStart == -1 {
				mapStart = i
			}
		} else

		// Check if the current index is the end of the map
		// Also Make sure the map start has been established
		if cache.data[i] == '}' && mapStart > 0 {

			// Check whether the query index is in between the map start
			// and the map end
			if valueIndex > mapStart && valueIndex < i+1 {

				// Append the data to the result array
				Result = append(Result, string(cache.data[mapStart+1:i]))

				// Reste the indexing variables
				mapStart, valueIndex = -1, -1
			}

		}
	}
	// Return the result
	return Result
}
