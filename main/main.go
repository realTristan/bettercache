package main

import (
	cache "github.com/realTristan/BetterCache"
)

func main() {
	var Cache *cache.Cache = cache.Init(-1)
	Cache.TestFullTextSearch(1, "my {name is }daniel!", "tristan", true)
}
