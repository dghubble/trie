package trie

// pathTrie is a trie of paths with string keys and generic type values.

// pathTrie is a trie of string keys and generic type values. By default
// PathTrie will segment keys by forward slashes with PathSegmenter
// (e.g. "/a/b/c" -> "/a", "/b", "/c"). A custom StringSegmenter may be
// used to customize how strings are segmented into nodes. A classic
// trie might segment keys by rune (i.e. unicode points).
type pathTrie[T any] struct {
	segmenter StringSegmenter // key segmenter, must not cause heap allocs
	value     *T
	children  map[string]*pathTrie[T]
}

// PathTrieOption is an optional configuration option for a path trie.
type PathTrieOption[T any] func(*pathTrie[T])

// WithSegmenter sets a non-default StringSegmenter on the path trie.
func WithSegmenter[T any](s StringSegmenter) PathTrieOption[T] {
	return func(trie *pathTrie[T]) { trie.segmenter = s }
}

// NewPathTrie allocates and returns a new path implementation of Trie.
func NewPathTrie[T any](opts ...PathTrieOption[T]) Trie[T] {
	trie := &pathTrie[T]{
		segmenter: PathSegmenter,
	}
	for _, opt := range opts {
		opt(trie)
	}
	return trie
}

// newPathTrieFromTrie returns new trie while preserving its config
func (trie *pathTrie[T]) newPathTrieFromTrie() *pathTrie[T] {
	return &pathTrie[T]{
		segmenter: trie.segmenter,
	}
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *pathTrie[T]) Get(key string) (T, bool) {
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		node = node.children[part]
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
func (trie *pathTrie[T]) Put(key string, value T) bool {
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		child := node.children[part]
		if child == nil {
			if node.children == nil {
				node.children = map[string]*pathTrie[T]{}
			}
			child = trie.newPathTrieFromTrie()
			node.children[part] = child
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
func (trie *pathTrie[T]) Delete(key string) bool {
	var path []nodeStr[T] // record ancestors to check later
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		path = append(path, nodeStr[T]{part: part, node: node})
		node = node.children[part]
		if node == nil {
			// node does not exist
			return false
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
func (trie *pathTrie[T]) Walk(walker WalkFunc[T]) error {
	return trie.walk("", walker)
}

// WalkPath iterates over each key/value in the path in trie from the root to
// the node at the given key, calling the given walker function for each
// key/value. If the walker function returns an error, the walk is aborted.
func (trie *pathTrie[T]) WalkPath(key string, walker WalkFunc[T]) error {
	// Get root value if one exists.
	if trie.value != nil {
		if err := walker("", *trie.value); err != nil {
			return err
		}
	}
	for part, i := trie.segmenter(key, 0); ; part, i = trie.segmenter(key, i) {
		if trie = trie.children[part]; trie == nil {
			return nil
		}
		if trie.value != nil {
			var k string
			if i == -1 {
				k = key
			} else {
				k = key[0:i]
			}
			if err := walker(k, *trie.value); err != nil {
				return err
			}
		}
		if i == -1 {
			break
		}
	}
	return nil
}

// PathTrie node and the part string key of the child the path descends into.
type nodeStr[T any] struct {
	node *pathTrie[T]
	part string
}

func (trie *pathTrie[T]) walk(key string, walker WalkFunc[T]) error {
	if trie.value != nil {
		if err := walker(key, *trie.value); err != nil {
			return err
		}
	}
	for part, child := range trie.children {
		if err := child.walk(key+part, walker); err != nil {
			return err
		}
	}
	return nil
}

func (trie *pathTrie[T]) isLeaf() bool {
	return len(trie.children) == 0
}
