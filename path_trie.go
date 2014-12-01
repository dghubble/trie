package trie

import (
	"strings"
)

// PathTrie is a trie of paths with string keys and interface{} values.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
type PathTrie struct {
	value    interface{}
	children map[string]*PathTrie
}

// New allocates and returns a new *PathTrie.
func NewPathTrie() *PathTrie {
	return &PathTrie{
		children: make(map[string]*PathTrie),
	}
}

// Part returns the next part of
func keySplit(path string, start int) (segment string, next int) {
	if len(path) == 0 {
		return path, -1
	}
	end := strings.IndexRune(path[start+1:], '/')
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *PathTrie) Get(key string) interface{} {
	node := trie
	for part, i := keySplit(key, 0); ; part, i = keySplit(key, i) {
		node = node.children[part]
		if node == nil {
			return nil
		}
		if i == -1 {
			break
		}
	}
	return node.value
}

// Put inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *PathTrie) Put(key string, value interface{}) bool {
	node := trie
	for part, i := keySplit(key, 0); ; part, i = keySplit(key, i) {
		child, _ := node.children[part]
		if child == nil {
			child = NewPathTrie()
			node.children[part] = child
		}
		node = child
		if i == -1 {
			break
		}
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = value
	return isNewVal
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the trie.
func (trie *PathTrie) Delete(key string) bool {
	path := make([]nodeStr, 0) // record ancestors to check later
	node := trie
	for part, i := keySplit(key, 0); ; part, i = keySplit(key, i) {
		path = append(path, nodeStr{part: part, node: node})
		node = node.children[part]
		if node == nil {
			// node does not exist
			return false
		}
		if i == -1 {
			break
		}
	}
	// delete the node value
	node.value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(path) - 1; i >= 0; i-- {
			parent := path[i].node
			part := path[i].part
			delete(parent.children, part)
			if parent.value != nil || !parent.isLeaf() {
				// parent has a value or has other children, stop
				break
			}
		}
	}
	return true // node (internal or not) existed and its value was nil'd
}

// Walk iterates over each key/value stored in the trie and calls the given
// walker function with the key and value. If the walker function returns
// an error, the walk is aborted.
// The traversal is depth first with no guaranteed order.
func (trie *PathTrie) Walk(walker WalkFunc) error {
	return trie.walk("", walker)
}

// PathTrie node and the part string key of the child the path descends into.
type nodeStr struct {
	node *PathTrie
	part string
}

func (trie *PathTrie) walk(key string, walker WalkFunc) error {
	if trie.value != nil {
		walker(key, trie.value)
	}
	for part, child := range trie.children {
		err := child.walk(key+part, walker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (trie *PathTrie) isLeaf() bool {
	return len(trie.children) == 0
}
