# Trie [![Build Status](https://travis-ci.org/dghubble/trie.png)](https://travis-ci.org/dghubble/trie) [![GoDoc](http://godoc.org/github.com/dghubble/trie?status.png)](http://godoc.org/github.com/dghubble/trie)

Trie implements several types of Tries (e.g. rune-wise). The implementations are optimized for ``Get`` performance and to allocate 0 bytes of heap memory (i.e. garbage) per Get.

The Tries do not synchronize access (not thread-safe). A typical use case is to perform ``Puts`` and ``Deletes`` upfront to populate the True, the perform ``Gets`` very quickly.

When Tries are chosen over maps, it is typically for their space efficiency. However, in situations where direct key lookup is not possible, tries can provide faster lookups (e.g. routers). 

## Install

    $ go get github.com/dghubble/trie

## Documentation

Read the full documentation: [https://godoc.org/github.com/dghubble/trie](https://godoc.org/github.com/dghubble/trie)

## Performance

These benchmarks write and read random string keys of byte length 30.

    RuneTrie (rune-wise)

    BenchmarkRuneTriePutStringKey  5000000         621 ns/op         1 B/op        0 allocs/op
    BenchmarkRuneTrieGetStringKey  5000000         611 ns/op         0 B/op        0 allocs/op

    BenchmarkRuneTriePutPathKey  5000000         687 ns/op         1 B/op        0 allocs/op
    BenchmarkRuneTrieGetPathKey  5000000         652 ns/op         0 B/op        0 allocs/op

## License

[MIT License](LICENSE)


