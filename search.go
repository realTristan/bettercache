package bettercache

import "strings"

// The TextSearch struct contains parameters for performing a search on the cache values.
//
// Parameters:
//   - Query: The string to search for in the cache values.
//   - Limit: The maximum number of search results to return. Set to -1 for no limit.
//   - StrictMode: Whether to perform a case-sensitive search.
type TextSearch struct {
	Query      string
	Limit      int
	StrictMode bool
}

// FullTextSearch performs a full text search using the cache values.
// It locks the cache mutex to prevent data overwriting before iterating over the cache data.
//
// The function uses the strings.Contains function to check whether the cache value contains the query.
// If the value contains the query, it will append the cache value to the result slice.
//
// If the user is not using strict mode, it will convert the cache data and the query to lowercase before
// checking whether the cache data contains the query.
//
// Parameters:
//   - TS: A pointer to a TextSearch struct that contains the query to search for, the maximum number of search results to return,
//     whether to perform a case-sensitive search, and a map of previous search queries.
//
// Returns:
//   - A slice of cache values that contain the provided query.
func (c *Cache) FullTextSearch(TS *TextSearch) []string {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// res -> The result slice containing the values
	// that contain the TextSearch.Query
	var res []string

	// Iterate over the cache data
	for i := 0; i < len(c.ft.data); i++ {

		// If the current c.ft.data index contains the
		// provided query
		if func() bool {
			// If the user is not using strict mode
			if !TS.StrictMode {
				// Convert the cache data and the query to lowercase
				// then return whether the cache data contains the query
				return strings.Contains(
					strings.ToLower(c.ft.data[i]),
					strings.ToLower(TS.Query))
			}

			// If the user is using strict mode, then just return
			// whether the cache data contains the query with no adjustments
			return strings.Contains(c.ft.data[i], TS.Query)
		}() {
			// Append value that contains the query to
			// the result slice
			res = append(res, c.ft.data[i])
		}
	}
	// Return the result slice
	return res
}
