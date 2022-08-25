# Better Cache ![Stars](https://img.shields.io/github/stars/realTristan/BetterCache?color=brightgreen) ![Watchers](https://img.shields.io/github/watchers/realTristan/BetterCache?label=Watchers)
![banner (5)](https://user-images.githubusercontent.com/75189508/186757681-6b7f97e8-ec37-448a-83cc-75106ed16309.png)

BetterCache, A lightweight, lightning fast caching system.

# About
- Better Cache is an ultra fast caching system that uses an array of bytes for storing data instead of the very common map caching system. Because the data is stored in an array of bytes instead of a map, it allows the user to achieve lightning fast full text search speeds. (within microseconds to milliseconds)
- Better Cache also only uses native golang modules which makes it fast, lightweight and easy to use.

# Benchmarks

<h3>Full Text Search</h3>
```
    Cache Size: 25 -> ~42µs
    Cache Size: 250 -> ~82µs
    Cache Size: 2500 -> ~220µs
    Cache Size: 25000 -> ~1.63ms
    Cache Size: 250000 -> ~16.87ms
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
