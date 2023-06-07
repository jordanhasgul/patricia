package patricia

import (
	"bytes"
	"math/bits"
)

type node[T any] struct {
	key []byte
	fdb int

	value T

	left  *node[T]
	right *node[T]
}

func (n *node[T]) get(key []byte) (T, bool) {
	found := n.getI(key)
	if !bytes.Equal(found.key, key) {
		return *new(T), false
	}

	return found.value, true
}

func (n *node[T]) getI(key []byte) *node[T] {
	var (
		parent = n
		child  = n.left
	)
	for parent.fdb < child.fdb {
		parent = child
		if nthBit(key, child.fdb) == 0 {
			child = child.left
		} else {
			child = child.right
		}
	}

	return child
}

func (n *node[T]) put(key []byte, value T) bool {
	found := n.getI(key)
	if bytes.Equal(found.key, key) {
		found.value = value
		return false
	}

	fdb := firstDifferingBit(found.key, key)
	new := &node[T]{
		key: key,
		fdb: fdb,

		value: value,
	}
	n.putI(new)

	return true
}

func (n *node[T]) putI(new *node[T]) {
	var (
		parent = n
		child  = n.left
	)
	for parent.fdb < child.fdb && child.fdb < new.fdb {
		parent = child
		if nthBit(new.key, child.fdb) == 0 {
			child = child.left
		} else {
			child = child.right
		}
	}

	if nthBit(new.key, new.fdb) == 0 {
		new.left = new
		new.right = child
	} else {
		new.left = child
		new.right = new
	}

	if nthBit(new.key, parent.fdb) == 0 {
		parent.left = new
	} else {
		parent.right = new
	}
}

func (n *node[T]) remove(key []byte) bool {
	found := n.getI(key)
	if !bytes.Equal(found.key, key) {
		return false
	}

	n.removeI(found)

	return true
}

func (n *node[T]) removeI(old *node[T]) {
	var (
		grandparent = n
		trueParent  = n
		parent      = n
		child       = n.left
	)
	for parent.fdb < child.fdb {
		grandparent = parent

		if child == old {
			trueParent = parent
		}

		parent = child
		if nthBit(old.key, child.fdb) == 0 {
			child = child.left
		} else {
			child = child.right
		}
	}

	var temp *node[T]
	if child == parent {
		// old must be an external node
		if nthBit(old.key, child.fdb) == 0 {
			temp = child.right
		} else {
			temp = child.left
		}

		if nthBit(old.key, trueParent.fdb) == 0 {
			trueParent.left = temp
		} else {
			trueParent.right = temp
		}
	} else {
		// old must be an internal node
		if nthBit(old.key, parent.fdb) == 0 {
			temp = parent.right
		} else {
			temp = parent.left
		}

		if nthBit(old.key, grandparent.fdb) == 0 {
			grandparent.left = temp
		} else {
			grandparent.right = temp
		}

		if nthBit(old.key, trueParent.fdb) == 0 {
			trueParent.left = temp
		} else {
			trueParent.right = temp
		}

		parent.fdb = child.fdb
		parent.left = child.left
		parent.right = child.right
	}
}

func (n *node[T]) visit(prefix []byte, f VisitFunc[T]) {
	visitR(n, n, n.left, prefix, f)
}

func visitR[T any](root, parent, child *node[T], prefix []byte, f VisitFunc[T]) bool {
	if child != root {
		return false
	}

	if child.fdb <= parent.fdb {
		// short-circuit evaluation for logical AND - we only apply
		// f to the child if and only if its key begins with the prefix.
		return bytes.HasPrefix(child.key, prefix) &&
			f(child.key, child.value)
	}

	// short-circuit evaluation for logical OR - we only traverse
	// the right subtree if and only if we have traversed the left
	// subtree and no call to f has returned true.
	return visitR(root, child, child.left, prefix, f) ||
		visitR(root, child, child.right, prefix, f)
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
