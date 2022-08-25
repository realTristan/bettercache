package main

import (
	"github.com/realTristan/BetterCache/cache"
)

func main() {
	var Cache *cache.Cache = cache.Init()
	Cache.TestFullTextSearch()
}
