package cache

// Import strings package
import (
	"strings"
)

/*

Notes:
Search how to store a type in a variable

then in the text search struct store it as SearchType
and when converting the interface{}, instead of .(string)
use .(SearchType)

*/

// The TextSearch struct contains five primary keys
/* Query: string { "The string to search for in the cache values" } 		*/
/* Limit: int { "The amount of search results. (Set to -1 for no limit)" }  */
/* StrictMode: bool { "Set to false to ignore caps in query comparisons" }  */
/* StorePreviousSearch: bool { "Set to true to keep previous query's" } 	*/
/* PreviousSearch: map[string][]string { "The Previous Searches" } 			*/
type TextSearch struct {
	Query               string
	Limit               int
	StrictMode          bool
	StorePreviousSearch bool
	PreviousSearch      map[string][]string
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
/* Parameters */
/* 	TS: *TextSearch = &TextSearch{
		Query               	string
		Limit               	int
		StrictMode          	bool
		StorePreviousSearch 	bool
		PreviousSearch      	map[string][]string
})
*/
//
// If you want to store the previous text search you made, you can set the
// StorePreviousSearch to true. This will set the key in the previous search
// to the provided TextSearch.Query and the value to the result slice.
//
/* >> Returns 			*/
/* res: []string	 	*/
func (cache *Cache) FullTextSearch(TS *TextSearch) []string {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// res -> The result slice containing the values
	// that contain the TextSearch.Query
	var res []string

	// Iterate over the cache data
	for i := 0; i < len(cache.fullTextData); i++ {

		// If the current cache.fullTextData index contains the
		// provided query
		if func() bool {
			// Make sure the current cache value was set
			// to true for full text search. If not, return false
			// if !strings.StartsWith(":FT(true):") {
			//	return false
			// }

			// If the user is not using strict mode
			if !TS.StrictMode {
				// Convert the cache data and the query to lowercase
				// then return whether the cache data contains the query
				return strings.Contains(
					strings.ToLower(cache.fullTextData[i]),
					strings.ToLower(TS.Query))
			}

			// If the user is using strict mode, then just return
			// whether the cache data contains the query with no adjustments
			return strings.Contains(cache.fullTextData[i], TS.Query)
		}() {
			// Append value that contains the query to
			// the result slice
			res = append(res, strings.Split(cache.fullTextData[i], ":")[1])
		}
	}
	// Add the result to the previous search
	// if the user set the previous search bool
	// to true.
	if TS.StorePreviousSearch {
		TS.PreviousSearch[TS.Query] = res
	}
	// Return the result slice
	return res
}
