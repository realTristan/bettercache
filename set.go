package bettercache

// Import fmt Package
import "fmt"

// The SetData struct has three primary keys
/* Key: interface{} { "The Cache Key" }									*/
/* Value: interface{} { "The Cache Value" }								*/
/* FullText: bool { "Whether to enable full text functions " }			*/
// WARNING
/* If FullText is set to true, it converts the Value to a string		*/
type SetData struct {
	Key      interface{}
	Value    interface{}
	FullText bool
}

// The Set function is used for setting a new value inside
// the cache data. The Set function locks the cache mutex to
// prevent data overwriting before checking whether the provided
// key already exists. If it does, it will call the Remove() function
// to remove that key from the cache.
//
// Once finished with the removal process, the function
// adds the value to the cache data slice, then adds the value's
// index to the cache indices map

// Sets a key to the provided value
/* Parameters: */
/* 	SD &SetData = *SetData{
		Key: interface{},
		Value: interface{},
		FullText: bool,
} */
func (c *Cache) Set(SD *SetData) {
	// Mutex locking
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If key exists
	if c.existsInMap(SD.Key) || c.existsInFullText(SD.Key) {
		// Remove the key from the cache
		c.remove(SD.Key)
	}

	// If the user set the AddToFullText to true
	if SD.FullText {
		// Set the key in the cache full text indices map
		// to the index the key value is at.
		c.fullTextIndices[SD.Key] = len(c.fullTextData)

		// Add the value into the cache data slice
		// as a modified string
		c.fullTextData = append(c.fullTextData,
			fmt.Sprintf("%s:%v", SD.Key, SD.Value))
	} else {
		// Set the key in the cache indices map to the
		// index the key value is at.
		c.mapData[SD.Key] = SD.Value
	}
	// Increase the current cache size
	c.currentSize++
}
