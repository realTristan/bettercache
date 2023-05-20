package bettercache

// Import the strings package
import "strings"

// The FullTextRemove function is used for removing cache keys whose values contain the provided query.
// The function iterates over the full text cache data and removes the keys whose values contain the query.
// The function returns a slice of removed keys.
//
// Parameters:
//   - query: A string representing the query to search for.
//   - amount: An integer representing the maximum number of keys to remove. If set to 0, all matching keys will be removed.
//
// Returns:
//   - A slice of removed keys.
func (c *Cache) FullTextRemove(query string, amount int) []string {
	// Mutex locking
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// res -> The result slice containing the values
	// that contain the TextSearch.Query
	var res []string

	// Iterate over the cache data
	for i := 0; i < len(c.ft.data); i++ {
		if len(res) >= amount && amount > 0 {
			return res
		} else if strings.Contains(c.ft.data[i], query) {
			// Append value that contains the query to
			// the result slice
			res = append(res, c.ft.data[i])

			// Remove the key
			c.remove(c.ft.data[i])
		}
	}

	// Return the slice containing all
	// of the removed values
	return res
}

// FullTextRemoveAll removes all cache keys that contain the provided query in their values.
// It utilizes the FullTextRemove function to remove the keys whose values contain the provided query.
//
// Parameters:
//   - query: The string to query for.
//
// Returns:
//   - A slice of removed keys.
func (c *Cache) FullTextRemoveAll(query string) []string {
	return c.FullTextRemove(query, -1)
}
