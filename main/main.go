package main

// Import packages
import (
	"fmt"
	"strings"
	"time"

	cache "github.com/realTristan/BetterCache"
)

/*

ask permission from teacher first though!

when have wifi, created a github repo with the class
code of the grade 12 computer science class I will be taking

then post anything from that class to that github repo

*/

// Main function
func main() {
	c := cache.Init(-1)
	c.TestFullTextSearch(10000, "1111111111111111111111111", "1", false)
	Map_VS_Slice()
}

func Map_VS_Slice() {
	fmt.Println("\nMap VS Slice Results:")
	st := time.Now()
	t1 := []interface{}{
		"1111", "1111", "1111", "1111", "1111", "1111", "1111", "1111", "1111", "1111", "1111",
	}
	for i := 0; i < len(t1); i++ {
		strings.Contains(t1[i].(string), "1")
	}
	fmt.Printf("Slice: %v\n", time.Since(st))

	st = time.Now()
	t2 := map[string]interface{}{
		"key1":  "1111",
		"key2":  "1111",
		"key3":  "1111",
		"key4":  "1111",
		"key5":  "1111",
		"key6":  "1111",
		"key7":  "1111",
		"key8":  "1111",
		"key9":  "1111",
		"key10": "1111",
		"key11": "1111",
	}
	for _, v := range t2 {
		strings.Contains(v.(string), "1")
	}
	fmt.Printf("Map: %v\n\n", time.Since(st))
}
