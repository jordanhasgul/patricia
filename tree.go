package patricia

import "bytes"

// Tree represents a generic Patricia tree i.e. a radix tree
// with radix r = 2. It is an associative data structure that
// maintains a collection of key-value pairs of type []byte
// and type T, respectively.
type Tree[T any] struct {
	root   *node[T]
	rooted bool

	size int
}

// New returns a pointer to a zero-valued [patricia.Tree].
func New[T any]() *Tree[T] {
	root := &node[T]{
		key: []byte{},
		fdb: -1,
	}
	root.left = root
	root.right = nil

	return &Tree[T]{root: root}
}

// Get returns the value asssociated with the key and a boolean
// indicating whether the key was present within the tree.
//
// Note: successive calls to Get are idempotent.
func (t *Tree[T]) Get(key []byte) (T, bool) {
	if bytes.Equal(t.root.key, key) {
		return t.root.value, t.rooted
	}

	return t.root.get(key)
}

// Put creates an association between the key and the value. If
// the key was present within the tree, its associated value is
// overwritten.
//
// Note: successive calls to Put are idempotent.
func (t *Tree[T]) Put(key []byte, value T) {
	if bytes.Equal(t.root.key, key) {
		t.root.value = value
		if !t.rooted {
			t.rooted = true
			t.size++
		}
		return
	}

	if t.root.put(key, value) {
		t.size++
	}
}

// Remove deletes any associations involving the key.
//
// Note: successive calls to Remove are idempotent.
func (t *Tree[T]) Remove(key []byte) {
	if bytes.Equal(t.root.key, key) {
		t.root.value = *new(T)
		if t.rooted {
			t.rooted = false
			t.size--
		}
		return
	}

	if t.root.remove(key) {
		t.size--
	}
}

type WalkFunc[T any] func([]byte, T) bool

// Walk performs an in-order traversal of the tree, applying
// f to those key-value pairs whose key begins with the prefix.
//
// The traversal terminates once f has returned true. Otherwise,
// it terminates once f has been applied to each key-value pair.
func (t *Tree[T]) Walk(prefix []byte, f WalkFunc[T]) {
	if bytes.Equal(t.root.key, prefix) {
		if t.rooted && f(t.root.key, t.root.value) {
			return
		}
	}

	t.root.walk(prefix, f)
}

// Size returns the number of key-value pairs present within
// the tree.
func (t *Tree[T]) Size() int {
	return t.size
}
