package main

import (
	cache "github.com/realTristan/BetterCache"
)

func main() {
	var Cache *cache.Cache = cache.Init(-1)

	Cache.TestFullTextSearch(2, "my {name is }tristan!!!!!!", "tristan", true)
}
