package patricia

import "bytes"

// Tree represents a generic Patricia tree i.e. a radix tree
// with radix r = 2. It is an associative data structure that
// maintains a collection of key-value pairs of type []byte
// and type T, respectively.
type Tree[T any] struct {
	root    *node[T]
	rootSet bool

	size int
}

// New returns a pointer to a zero-valued [patricia.Tree].
func New[T any]() *Tree[T] {
	root := &node[T]{
		key: []byte{},
		fdb: -1,

		value: *new(T),
	}
	root.left = root
	root.right = nil

	return &Tree[T]{
		root:    root,
		rootSet: false,

		size: 0,
	}
}

// Get returns the value asssociated with the key and a boolean
// indicating whether the key was present within the tree.
//
// Note: calls to Get are idempotent.
func (t *Tree[T]) Get(key []byte) (T, bool) {
	if bytes.Equal(t.root.key, key) {
		return t.root.value, t.rootSet
	}

	return t.root.get(key)
}

// Put creates an association between the key and the value. If
// the key was present within the tree, its associated value is
// overwritten.
//
// Note: calls to Put are idempotent.
func (t *Tree[T]) Put(key []byte, value T) {
	if bytes.Equal(t.root.key, key) {
		t.root.value = value
		if !t.rootSet {
			t.rootSet = true
			t.size++
		}
		return
	}

	if t.root.put(key, value) {
		t.size++
	}
}

// Remove deletes any association involving the key and a value.
//
// Note: calls to Remove are idempotent.
func (t *Tree[T]) Remove(key []byte) {
	if bytes.Equal(t.root.key, key) {
		t.root.value = *new(T)
		if t.rootSet {
			t.rootSet = false
			t.size--
		}
		return
	}

	if t.root.remove(key) {
		t.size--
	}
}

type VisitFunc[T any] func([]byte, T) bool

// Visit traverses the tree in-order and applies f to any key-value
// pairs whose key begins with the prefix.
//
// The traversal terminates once f has returned true. Otherwise, it
// terminates once every key-value pair has been visited.
func (t *Tree[T]) Visit(prefix []byte, f VisitFunc[T]) {
	if t.rootSet && bytes.HasPrefix(t.root.key, prefix) {
		if f(t.root.key, t.root.value) {
			return
		}
	}

	t.root.visit(prefix, f)
}

// Size returns the number of associations within the tree.
func (t *Tree[T]) Size() int {
	return t.size
}
