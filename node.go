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

func firstDifferingBit(buf1, buf2 []byte) int {
	if len(buf1) == 0 && len(buf2) != 0 {
		return bits.LeadingZeros8(buf2[0])
	}

	if len(buf1) != 0 && len(buf2) == 0 {
		return bits.LeadingZeros8(buf1[0])
	}

	n := len(buf1)
	if len(buf2) > n {
		n = len(buf2)
	}

	newBuf1 := make([]byte, n)
	copy(newBuf1, buf1)

	newBuf2 := make([]byte, n)
	copy(newBuf2, buf2)

	var differingByte int
	for idx := 0; idx < n; idx++ {
		if newBuf1[idx] != newBuf2[idx] {
			differingByte = idx
			break
		}
	}

	differingBit := newBuf1[differingByte] ^ newBuf2[differingByte]
	return 8*differingByte + bits.LeadingZeros8(differingBit)
}
