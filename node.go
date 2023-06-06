package patricia

import "math/bits"

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

func nthBit(buf []byte, n int) byte {
	if n < 0 || n >= 8*len(buf) {
		return 0
	}

	msb := (buf[n/8] << (n % 8)) & 128
	return bits.RotateLeft8(msb, 1)
}
