package trie

import (
	"testing"
)

// test splitting /path/keys/ into parts (e.g. /path, /keys, /)
func TestPathSegmenter(t *testing.T) {
	cases := []struct {
		key     string
		parts   []string
		indices []int // indexes to use as next start, in order
	}{
		{"", []string{""}, []int{-1}},
		{"/", []string{"/"}, []int{-1}},
		{"static_file", []string{"static_file"}, []int{-1}},
		{"/users/scott", []string{"/users", "/scott"}, []int{6, -1}},
		{"users/scott", []string{"users", "/scott"}, []int{5, -1}},
		{"/users/ramona/", []string{"/users", "/ramona", "/"}, []int{6, 13, -1}},
		{"users/ramona/", []string{"users", "/ramona", "/"}, []int{5, 12, -1}},
		{"//", []string{"/", "/"}, []int{1, -1}},
		{"/a/b/c", []string{"/a", "/b", "/c"}, []int{2, 4, -1}},
	}

	for _, c := range cases {
		partNum := 0
		for prefix, i := PathSegmenter(c.key, 0); ; prefix, i = PathSegmenter(c.key, i) {
			if prefix != c.parts[partNum] {
				t.Errorf("expected part %d of key '%s' to be '%s', got '%s'", partNum, c.key, c.parts[partNum], prefix)
			}
			if i != c.indices[partNum] {
				t.Errorf("in iteration %d, expected next index of key '%s' to be '%d', got '%d'", partNum, c.key, c.indices[partNum], i)
			}
			partNum++
			if i == -1 {
				break
			}
		}
		if partNum != len(c.parts) {
			t.Errorf("expected '%s' to have %d parts, got %d", c.key, len(c.parts), partNum)
		}
	}
}
