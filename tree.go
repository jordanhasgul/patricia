package patricia

type Tree[T any] struct {
	root    *node[T]
	rootSet bool

	size int
}

func New[T any]() *Tree[T] {
	return nil
}

func (t *Tree[T]) Get(key []byte) (T, bool) {
	return *new(T), false
}

func (t *Tree[T]) Put(key []byte, value T) {

}

func (t *Tree[T]) Remove(key []byte) {

}

func (t *Tree[T]) Size() int {
	return 0
}
