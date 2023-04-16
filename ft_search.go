package bettercache

// Import strings package
import "strings"

// The TextSearch struct contains five primary keys
/* Query: string { "The string to search for in the cache values" } 		*/
/* Limit: int { "The amount of search results. (Set to -1 for no limit)" }  */
/* StrictMode: bool { "Set to false to ignore caps in query comparisons" }  */
/* PreviousQueries: map[string][]string { "The Previous Searches" } 			*/
type TextSearch struct {
	Query           string
	Limit           int
	StrictMode      bool
	PreviousQueries map[string][]string
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
/* Parameters: */
/* 	TS: *TextSearch = &TextSearch{
		Query               	string
		Limit               	int
		StrictMode          	bool
		PreviousQueries      		map[string][]string
})
*/
//
// If you want to store the previous text search you made, you can set the
// PreviousQueries map. This will set the key in the previous search
// to the provided TextSearch.Query and the value to the result slice.
// It is suggested to use a limited sized map.
//
/* >> Returns 			*/
/* res: []string	 	*/
func (c *Cache) FullTextSearch(TS *TextSearch) []string {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// res -> The result slice containing the values
	// that contain the TextSearch.Query
	var res []string

	// Iterate over the cache data
	for i := 0; i < len(c.fullTextData); i++ {

		// If the current c.fullTextData index contains the
		// provided query
		if func() bool {
			// If the user is not using strict mode
			if !TS.StrictMode {
				// Convert the cache data and the query to lowercase
				// then return whether the cache data contains the query
				return strings.Contains(
					strings.ToLower(c.fullTextData[i]),
					strings.ToLower(TS.Query))
			}

			// If the user is using strict mode, then just return
			// whether the cache data contains the query with no adjustments
			return strings.Contains(c.fullTextData[i], TS.Query)
		}() {
			// Append value that contains the query to
			// the result slice
			res = append(res, strings.Split(c.fullTextData[i], ":")[1])
		}
	}
	// Add the result to the previous search
	// if the user set the previous search bool
	// to true.
	if TS.PreviousQueries != nil {
		TS.PreviousQueries[TS.Query] = res
	}
	// Return the result slice
	return res
}

// The GetPreviousQuery function is used to return the
// slice of values for a previous query
//
/* Paramters */
/* query: string { "The Previous Query" } */
//
/* Returns */
/* results: []string { "The Query Results" } */
func (TS *TextSearch) GetPreviousQuery(query string) []string {
	return TS.PreviousQueries[query]
}

// The GetPreviousQuery function is used to delete a
// previous query from the PreviousQueries map
//
/* Paramters */
/* query: string { "The Previous Query" } */
func (TS *TextSearch) DeletePreviousQuery(query string) {
	delete(TS.PreviousQueries, query)
}

// The ClearPreviousQueries function is used to reset
// the previous queries map
//
/* Paramters */
/* size: int { "The Size of PreviousQueries Map (Set to -1 for no limit)" } */
func (TS *TextSearch) ClearPreviousQueries(size int) {
	TS.PreviousQueries = map[string][]string{}
}
