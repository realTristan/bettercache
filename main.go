package main

// Main function
func main() {
	// Init the cache
	var Cache *Cache = Init(100) // 100 -> Size in bytes

	// Test the full text search
	Cache.TestFullTextSearch()
}
