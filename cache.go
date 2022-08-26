package cache

// Import modules
import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// The Cache struct contains three primary keys
/* Data: []byte -> The Cache Data in Bytes						 	 */
/* Mutex: *sync.Mutex -> Used for locking/unlocking the data 	 	 */
/* MaxSize int -> The cache max size 								 */
type Cache struct {
	data    []byte
	mutex   *sync.RWMutex
	maxSize int
}

// The initData() function returns a the cache
// byte slice depending on the provided size
func initData(size int) []byte {
	// Limited Size
	if size > 0 {
		var data = make([]byte, size+1)
		data[0] = '*'
		return data

	}
	// Unlimited size
	return []byte{'*'}
}

// The Init() function creates the Cache
// object depending on what was entered for
// the size of the cache
func Init(size int) *Cache {
	// Return the new cache
	return &Cache{
		data:    initData(size),
		mutex:   &sync.RWMutex{},
		maxSize: size,
	}
}

// The Exists() function returns whether the
// provided key exists in the cache
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) Exists(key string) bool {
	return len(cache.Get(key)) > 0
}

// The GetByteSize() function returns the current size of the
// cache bytes and the cache maximum size
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) GetByteSize() (int, int) {
	return len(cache.data), cache.maxSize
}

// The Expire() function removes the provided key
// from the cache after the given time
//
// Suggested to run this function in your
// own monitored goroutine
func (cache *Cache) Expire(key string, _time time.Duration) {
	time.Sleep(_time)
	cache.Remove(key)
}

// The Flush() function resets the cache
// data. Make sure to use this function when
// clearing the cache!
func (cache *Cache) Flush() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Reset the cache data
	cache.data = []byte{'*'}
}

// The ShowBytes() function returns the cache bytes
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) ShowBytes() []byte {
	return cache.data
}

// The Show() function returns the cache as a string
//
// No read lock/unlock because this function isn't
// as heavy as the ones that do utilize the read locks
func (cache *Cache) Show() string {
	return string(cache.data)
}

// The Set() function sets the value for the
// provided key inside the cache.
//
// Example: {"key1": "my name is tristan!"},
//
// Returns the removed value of the previously
// defined key
func (cache *Cache) Set(key string, data string) string {
	// Define Variables
	var (
		// removedValue -> The previous value removed
		removedValue string = ""
		// keyBytes -> The modified key in a bytes slice
		keyBytes []byte = []byte(fmt.Sprintf(`%s:%d{`, key, len(data)))
	)

	// I put this in it's own function so that it can
	// use defer for unlocking the mutex.
	// This is best because if there's any errors, they
	// won't prevent the mutex from unlocking
	if func() bool {
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		// Check if the key already exists
		return bytes.Contains(cache.data, []byte(key))
	}() {
		// Remove the previous key
		removedValue = cache.Remove(key)
	}
	// Lock/Unlock the mutex
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Set the byte cache value
	cache.data = append(
		cache.data, append(keyBytes, append([]byte(data), '}')...)...)

	// Return the removed value
	return removedValue
}
