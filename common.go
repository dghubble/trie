package trie

import (
	"strings"
)

// A WalkFunc defines some action to take on the given key and value during
// a Trie Walk. Returning a non-nil error will terminate the Walk.
type WalkFunc func(key string, value interface{}) error

// A Segmenter function takes a string key with a starting index and returns
// the first segment after the start and the ending index. When the end is
// reached, the returned nextIndex should be -1.
// Implementations should NOT allocate heap memory as Trie Segmentors are
// called upon Gets. See PathSegmentor.
type StringSegmenter func(key string, start int) (segment string, nextIndex int)

// PathSegmenter segments string key paths by slash separators. For example,
// "/a/b/c" -> ("/a", 2), ("/b", 4), ("/c", -1) in successive calls. It does
// not allocate any heap memory.
func PathSegmenter(path string, start int) (segment string, next int) {
	if len(path) == 0 {
		return path, -1
	}
	end := strings.IndexRune(path[start+1:], '/')
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1
}
