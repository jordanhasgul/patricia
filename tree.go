package patricia

import "bytes"

type Tree[T any] struct {
	root    *node[T]
	rootSet bool

	size int
}

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

func (t *Tree[T]) Get(key []byte) (T, bool) {
	if bytes.Equal(t.root.key, key) {
		return t.root.value, t.rootSet
	}

	return t.root.get(key)
}

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

func (t *Tree[T]) Size() int {
	return t.size
}
