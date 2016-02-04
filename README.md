# Trie [![Build Status](https://travis-ci.org/dghubble/trie.png)](https://travis-ci.org/dghubble/trie) [![Coverage](https://gocover.io/_badge/github.com/dghubble/trie)](https://gocover.io/github.com/dghubble/trie) [![GoDoc](https://godoc.org/github.com/dghubble/trie?status.png)](https://godoc.org/github.com/dghubble/trie)

Trie implements several types of Tries (e.g. rune-wise, path-wise). The implementations are optimized for ``Get`` performance and to allocate 0 bytes of heap memory (i.e. garbage) per Get.

The Tries do not synchronize access (not thread-safe). A typical use case is to perform ``Puts`` and ``Deletes`` upfront to populate the Trie, then perform ``Gets`` very quickly.

When Tries are chosen over maps, it is typically for their space efficiency. However, in situations where direct key lookup is not possible (e.g. routers), tries can provide faster lookups and avoid key iteration. 

## Install

    $ go get github.com/dghubble/trie

## Documentation

Read [Godoc](https://godoc.org/github.com/dghubble/trie)

## Performance

RuneTrie is a typical Trie which segments strings rune-wise (i.e. by unicode code point). These benchmarks perform Puts and Gets of random string keys that are 30 bytes long and of random '/' separated paths that have 3 parts and are 30 bytes long (longer if you count the '/' seps).

    BenchmarkRuneTriePutStringKey    2000000      653 ns/op      2 B/op     0 allocs/op
    BenchmarkRuneTrieGetStringKey    5000000      616 ns/op      0 B/op     0 allocs/op
    BenchmarkRuneTriePutPathKey      5000000      704 ns/op      1 B/op     0 allocs/op
    BenchmarkRuneTrieGetPathKey      5000000      682 ns/op      0 B/op     0 allocs/op

PathTrie segments strings by forward slash separators which can boost performance
for some use cases. These benchmarks perform Puts and Gets of random string keys that are 30 bytes long and of random '/' separated paths that have 3 parts and are 30 bytes long (longer if you count the '/' seps).

    BenchmarkPathTriePutStringKey   20000000     94.2 ns/op      0 B/op     0 allocs/op
    BenchmarkPathTrieGetStringKey   20000000     93.5 ns/op      0 B/op     0 allocs/op
    BenchmarkPathTriePutPathKey     20000000      113 ns/op      0 B/op     0 allocs/op
    BenchmarkPathTrieGetPathKey     20000000      108 ns/op      0 B/op     0 allocs/op

Note that for random string Puts and Gets, the PathTrie is effectively a map as every node is a direct child of the root (except for strings that happen to have a slash).

This benchmark measures the performance of the PathSegmenter alone. It is used to segment random paths that have 3 '/' separated parts and are 30 bytes long.

    BenchmarkPathSegmenter          50000000     58.8 ns/op      0 B/op     0 allocs/op

## License

[MIT License](LICENSE)


