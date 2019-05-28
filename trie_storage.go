package trie

import (
	"encoding/gob"
	"fmt"
	"io"
)

type trieStorage struct {
	// not storable on disk
	// Segmenter StringSegmenter
	Value    interface{}
	Children map[string]*trieStorage
}

// Save stores the trie to file in a gob format.
// does not store the Segmenter function
func (trie *PathTrie) Save(out io.Writer) error {
	ts := trie.toTrieStorage()
	return gob.NewEncoder(out).Encode(&ts)
}

// Save stores a rune-trie
func (trie *RuneTrie) Save(out io.Writer) error {
	ts := trie.toTrieStorage()
	return gob.NewEncoder(out).Encode(ts)
}

// LoadPathTrie reads a PathTrie from file (or io.Reader)
func LoadPathTrie(in io.Reader) (*PathTrie, error) {
	var ts trieStorage
	err := gob.NewDecoder(in).Decode(&ts)
	if err != nil {
		return nil, err
	}
	return ts.toPathTrie(), nil
}

// LoadRuneTrie reads a RuneTrie from file (or any io.Reader)
func LoadRuneTrie(in io.Reader) (*RuneTrie, error) {
	var ts trieStorage
	err := gob.NewDecoder(in).Decode(&ts)
	if err != nil {
		return nil, err
	}
	return ts.toRuneTrie()
}

func (ts *trieStorage) toPathTrie() *PathTrie {
	trie := &PathTrie{segmenter: PathSegmenter, value: ts.Value, children: make(map[string]*PathTrie)}
	for k, v := range ts.Children {
		trie.children[k] = v.toPathTrie()
	}
	return trie
}

func (ts *trieStorage) toRuneTrie() (*RuneTrie, error) {
	trie := &RuneTrie{value: ts.Value, children: make(map[rune]*RuneTrie)}
	for k, v := range ts.Children {
		rs := []rune(k)
		if len(rs) != 1 {
			return nil, fmt.Errorf("not a rune trie")
		}
		child, err := v.toRuneTrie()
		if err != nil {
			return nil, err
		}
		trie.children[rs[0]] = child
	}
	return trie, nil
}

func (trie *PathTrie) toTrieStorage() *trieStorage {
	ts := &trieStorage{Value: trie.value, Children: make(map[string]*trieStorage)}
	for k, v := range trie.children {
		ts.Children[k] = v.toTrieStorage()
	}
	return ts
}

func (trie *RuneTrie) toTrieStorage() *trieStorage {
	ts := &trieStorage{Value: trie.value, Children: make(map[string]*trieStorage)}
	for k, v := range trie.children {
		ts.Children[string(k)] = v.toTrieStorage()
	}
	return ts
}
