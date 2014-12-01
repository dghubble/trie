package trie

import (
	"testing"
)

// test splitting /path/keys/ into parts (e.g. /path, /keys, /)
func TestKeySpliter(t *testing.T) {
	cases := []struct {
		key   string
		parts []string
	}{
		{"", []string{""}},
		{"/", []string{"/"}},
		{"static_file", []string{"static_file"}},
		{"/users/scott", []string{"/users", "/scott"}},
		{"users/scott", []string{"users", "/scott"}},
		{"/users/ramona/", []string{"/users", "/ramona", "/"}},
		{"users/ramona/", []string{"users", "/ramona", "/"}},
		{"//", []string{"/", "/"}},
	}

	for _, c := range cases {
		partNum := 0
		for prefix, i := keySplit(c.key, 0); ; prefix, i = keySplit(c.key, i) {
			if prefix != c.parts[partNum] {
				t.Errorf("expected part %d of key '%s' to be '%s', got '%s'", partNum, c.key, c.parts[partNum], prefix)
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
