# Better Cache ![Stars](https://img.shields.io/github/stars/realTristan/bettercache?color=brightgreen) ![Watchers](https://img.shields.io/github/watchers/realTristan/bettercache?label=Watchers)
![banner (5)](https://user-images.githubusercontent.com/75189508/186757681-6b7f97e8-ec37-448a-83cc-75106ed16309.png)

# Install
`go get -u github.com/realTristan/bettercache`

# Example Usage
```go
package main

// Import Packages
import (
    "fmt"
    bc "github.com/realTristan/bettercache"
)

func main() {
    // Initialize the cache
    var c *bc.Cache = bc.Init(-1) // -1 (no size limit)

    // Add key1 to the cache
    c.Set("key1", "value1", true)

    // Get key from the cache
    var data string = c.Get("key1")
    fmt.Println(data)

    // Full Text Search for the key's contents
    var res []string = c.FullTextSearch(&bc.TextSearch{
        Limit:      -1,                 // No limit
        Query:      []byte("value"),    // Search for "value"
        StrictMode: false,              // Ignore CAPS
    })
    fmt.Println(res)

    // Remove key1 from the cache
    var removedKey string = c.Remove("key1")
    fmt.Println(removedKey)
}
```

# License
MIT License

Copyright (c) 2022 Tristan Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
