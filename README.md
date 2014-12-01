# Trie [![Build Status](https://travis-ci.org/dghubble/trie.png)](https://travis-ci.org/dghubble/trie) [![GoDoc](http://godoc.org/github.com/dghubble/trie?status.png)](http://godoc.org/github.com/dghubble/trie)

Trie implements several types of Tries (e.g. rune-wise, path-wise). The implementations are optimized for ``Get`` performance and to allocate 0 bytes of heap memory (i.e. garbage) per Get.

The Tries do not synchronize access (not thread-safe). A typical use case is to perform ``Puts`` and ``Deletes`` upfront to populate the Trie, then perform ``Gets`` very quickly.

When Tries are chosen over maps, it is typically for their space efficiency. However, in situations where direct key lookup is not possible (e.g. routers), tries can provide faster lookups and avoid key iteration. 

## Install

    $ go get github.com/dghubble/trie

## Documentation

Read the full documentation: [https://godoc.org/github.com/dghubble/trie](https://godoc.org/github.com/dghubble/trie)

## Performance

These benchmarks Put and Get random string keys (30 bytes long).

    BenchmarkRuneTriePutStringKey  5000000    613 ns/op    1 B/op   0 allocs/op
    BenchmarkRuneTrieGetStringKey  5000000    623 ns/op    0 B/op   0 allocs/op
    BenchmarkPathTriePutStringKey  20000000   92.0 ns/op   0 B/op   0 allocs/op
    BenchmarkPathTrieGetStringKey  20000000   96.3 ns/op   0 B/op   0 allocs/op

Note that for random string keys without '/' separators, a PathTrie is effectively a map as every prefix is a direct child of the root.

Putting and Getting paths with 3 slash separated parts, where each part is a string of 10 random bytes.

    BenchmarkRuneTriePutPathKey    5000000    679 ns/op    1 B/op   0 allocs/op
    BenchmarkRuneTrieGetPathKey    5000000    674 ns/op    0 B/op   0 allocs/op    
    BenchmarkPathTriePutPathKey    20000000   111 ns/op    0 B/op   0 allocs/op
    BenchmarkPathTrieGetPathKey    20000000   109 ns/op    0 B/op   0 allocs/op

Benchmark for the path key splitter used in the PathTrie.
    
    BenchmarkKeySplitter           50000000   58.7 ns/op   0 B/op   0 allocs/op

## License

[MIT License](LICENSE)


