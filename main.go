package main

// Main function
func main() {
	// Init the cache
	var Cache *Cache = Init(-1) // -1 -> Unlimited Size

	// Test the full text search
	Cache.TestFullTextSearch()
}
