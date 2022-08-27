package cache

// Import the strings package
import (
	"strings"
)

// The TextRemove struct has two primary keys
/* Query: string { "The string to find within values" } 				*/
/* Amount: int { "The amount of keys to remove (Set to -1 for all)" }	*/
type TextRemove struct {
	Query  string
	Amount int
}

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
/* Parameters */
/* 	TS: *TextRemove = &TextRemove{
		Query               	string
		Amount               	int
})*/
//
// If you want to remnove all the values, either call the FullTextRemoveAll()
// function or set the TextRemove.Amount to -1
//
/* >> Returns */
/* res: []string */
func (cache *Cache) FullTextRemove(TR *TextRemove) []string {
	// Mutex locking
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	// res -> The result slice containing the values
	// that contain the TextSearch.Query
	var res []string

	// Iterate over the cache data
	for i := 0; i < len(cache.fullTextData); i++ {
		if len(res) >= TR.Amount && TR.Amount > 0 {
			return res
		} else

		// If the current cache.fullTextData index contains the
		// provided query
		if strings.Contains(cache.fullTextData[i], TR.Query) {
			// Split the cache value by ':' to bypass
			// the {key_name}:
			var key string = strings.Split(cache.fullTextData[i], ":")[0]
			// Append value that contains the query to
			// the result slice
			res = append(res, cache.fullTextData[i][len(key)+1:])

			// I decided to put this inside a function so that
			// even if there's any errors in the Remove function,
			// the mutex will still relock once the function returns
			func() {
				// Unlock the mutex so the remove function can
				// remove the key from the cache
				cache.mutex.RUnlock()
				// Then re-lock the mutex once the key has
				// been removed
				defer cache.mutex.RLock()

				// Remove the key from the cache
				cache.Remove(key)
			}()
		}
		//}
	}
	// Return the slice containing all
	// of the removed values
	return res
}

// The Full Text Remove All function utilizes the Full Text Remove function
// to remove the keys whos values contain the provided query.

// Removes all cache keys that contain the provided query in their values
/* Paramters */
/*	query: string { "The string to query for" } */
func (cache *Cache) FullTextRemoveAll(query string) []string {
	return cache.FullTextRemove(&TextRemove{
		Query:  query,
		Amount: -1,
	})
}
