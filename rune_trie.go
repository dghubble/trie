package trie

// runeTrie is a trie of runes with string keys and generic type values.
type runeTrie[T any] struct {
	value    *T
	children map[rune]*runeTrie[T]
}

// NewRuneTrie allocates and returns a new rune implementation of Trie.
func NewRuneTrie[T any]() Trie[T] {
	return new(runeTrie[T])
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *runeTrie[T]) Get(key string) (T, bool) {
	node := trie
	for _, r := range key {
		node = node.children[r]
		if node == nil {
			return zeroValueOfT[T](), false
		}
	}
	if node.value == nil {
		return zeroValueOfT[T](), false
	}
	return *node.value, true
}

// Put inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *runeTrie[T]) Put(key string, value T) bool {
	node := trie
	for _, r := range key {
		child := node.children[r]
		if child == nil {
			if node.children == nil {
				node.children = map[rune]*runeTrie[T]{}
			}
			child = new(runeTrie[T])
			node.children[r] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = &value
	return isNewVal
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the trie.
func (trie *runeTrie[T]) Delete(key string) bool {
	path := make([]nodeRune[T], len(key)) // record ancestors to check later
	node := trie
	for i, r := range key {
		path[i] = nodeRune[T]{r: r, node: node}
		node = node.children[r]
		if node == nil {
			// node does not exist
			return false
		}
	}
	// delete the node value
	node.value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor
	// path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(key) - 1; i >= 0; i-- {
			if path[i].node == nil {
				continue
			}
			parent := path[i].node
			r := path[i].r
			delete(parent.children, r)
			if !parent.isLeaf() {
				// parent has other children, stop
				break
			}
			parent.children = nil
			if parent.value != nil {
				// parent has a value, stop
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
func (trie *runeTrie[T]) Walk(walker WalkFunc[T]) error {
	return trie.walk("", walker)
}

// WalkPath iterates over each key/value in the path in trie from the root to
// the node at the given key, calling the given walker function for each
// key/value. If the walker function returns an error, the walk is aborted.
func (trie *runeTrie[T]) WalkPath(key string, walker WalkFunc[T]) error {
	// Get root value if one exists.
	if trie.value != nil {
		if err := walker("", *trie.value); err != nil {
			return err
		}
	}

	for i, r := range key {
		if trie = trie.children[r]; trie == nil {
			return nil
		}
		if trie.value != nil {
			if err := walker(string(key[0:i+1]), *trie.value); err != nil {
				return err
			}
		}
	}
	return nil
}

// RuneTrie node and the rune key of the child the path descends into.
type nodeRune[T any] struct {
	node *runeTrie[T]
	r    rune
}

func (trie *runeTrie[T]) walk(key string, walker WalkFunc[T]) error {
	if trie.value != nil {
		if err := walker(key, *trie.value); err != nil {
			return err
		}
	}
	for r, child := range trie.children {
		if err := child.walk(key+string(r), walker); err != nil {
			return err
		}
	}
	return nil
}

func (trie *runeTrie[T]) isLeaf() bool {
	return len(trie.children) == 0
}
