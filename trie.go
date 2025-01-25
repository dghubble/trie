package trie

// Trie exposes the Trie structure capabilities.
type Trie[T any] interface {
	Get(key string) (T, bool)
	Put(key string, value T) bool
	Delete(key string) bool
	Walk(walker WalkFunc[T]) error
	WalkPath(key string, walker WalkFunc[T]) error
}
