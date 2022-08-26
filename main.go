package main

// Main function
func main() {
	// Init the cache
	var Cache *Cache = Init(-1) // -1 -> Unlimited Size

	// Test the full text search
	Cache.TestFullTextSearch()
}

/*

for the set and remove functions
make it iterate through the cache byte
and do what u did with the get function.
then take the current index - len(key)
and set that as the startIndex
then iterate to loop for the closing bracket (})
then remove all data from the startIndex to that closing
bracket.

Once done, append the bytes of the given value (map[string]string)
to the end of the cache bytes

this will solve the longer times because of the json marshal/unmarshal
*/
