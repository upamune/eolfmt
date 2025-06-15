// Package processor handles file processing and newline fixing.
package processor

import "sync"

// Buffer pool to reduce small allocations
var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 512)
		return &b
	},
}

func getBuf() *[]byte {
	return bufPool.Get().(*[]byte)
}

func putBuf(buf *[]byte) {
	bufPool.Put(buf)
}
