package copy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestMemoryCopy_Int(t *testing.T) {
	src := 123
	dst := 456
	MemoryCopy(unsafe.Pointer(&dst), unsafe.Pointer(&src), int(unsafe.Sizeof(src)))
	assert.Equal(t, 123, dst)
}

func TestMemoryCopy_Array(t *testing.T) {
	src := [...]byte{1, 2, 3, 4, 5}
	dst := [...]byte{0, 0, 0, 0, 0, 0}
	MemoryCopy(unsafe.Pointer(&dst), unsafe.Pointer(&src), int(unsafe.Sizeof(byte(1)))*3)
	assert.Equal(t, [...]byte{1, 2, 3, 0, 0, 0}, dst)
}
