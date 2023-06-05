package patricia

type node[T any] struct {
	key []byte
	fdb int

	value T

	left  *node[T]
	right *node[T]
}

func (n *node[T]) get(key []byte) (T, bool) {
	return *new(T), false
}

func (n *node[T]) put(key []byte, value T) bool {
	return false
}

func (n *node[T]) remove(key []byte) bool {
	return false
}
