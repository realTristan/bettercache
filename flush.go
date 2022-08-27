package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func (cache *Cache) flushToFile(path string) {
	// If the BetterCache file doesn't exist
	if _, err := os.Stat(path); err != nil {
		// Return the function
		return
	}

	// Create the result map
	var result []byte = []byte{}

	// Make sure the full text indices is greater than
	// zero. Having a check prevents potential bugs.
	if len(cache.fullTextIndices) > 0 {
		// Iterate over the full text indices
		for _, i := range cache.fullTextIndices {

			// Append the full text key and value to the
			// result byte array
			result = append(result,
				[]byte(fmt.Sprintf("$FULLTEXT:%s\n", cache.fullTextData[i]))...)
		}
	}

	// Make sure the cache map data is greater than
	// zero. Having a check prevents potential bugs.
	if len(cache.mapData) > 0 {

		// Iterate over the cache map data
		for k, v := range cache.mapData {

			// Append the cache map key and value to the
			// result byte array
			result = append(result,
				[]byte(fmt.Sprintf("$CACHE:%v:%s\n", k, v))...)
		}
	}

	// Write the result byte array to the BetterCache file
	os.WriteFile(path, result, 0644)
}

func (cache *Cache) handleCacheFileLine(path string, line string, lineCount int) {
	cache.Set(&SetData{
		Key:      "",
		Value:    "",
		FullText: strings.HasPrefix(line, "$FULLTEXT"),
	})
}

func (cache *Cache) ReadFlush(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Scan the file's lines
	var scanner *bufio.Scanner = bufio.NewScanner(file)
	var lineCount int = 0
	for scanner.Scan() {
		cache.handleCacheFileLine(path, scanner.Text(), lineCount)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

/*
use a json file instead



*/

// Flush the data to the BetterCache file
func (cache *Cache) Flush(path string) {
	// Mutex locking
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Flush cache to BetterCache file
	cache.flushToFile(path)
}
