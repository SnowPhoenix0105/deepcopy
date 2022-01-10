package memcpy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestMemoryCopy_Int(t *testing.T) {
	src := 123
	dst := 456
	Do(uintptr(unsafe.Pointer(&dst)), uintptr(unsafe.Pointer(&src)), unsafe.Sizeof(src))
	assert.Equal(t, 123, dst)
}

func TestMemoryCopy_Array3(t *testing.T) {
	src := [...]byte{1, 2, 3, 4, 5}
	dst := [...]byte{0, 0, 0, 0, 0, 0}
	Do(uintptr(unsafe.Pointer(&dst)), uintptr(unsafe.Pointer(&src)), unsafe.Sizeof(byte(1))*3)
	assert.Equal(t, [...]byte{1, 2, 3, 0, 0, 0}, dst)
}

func TestMemoryCopy_Array8(t *testing.T) {
	src := [11]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	dst := [11]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Do(uintptr(unsafe.Pointer(&dst))+2, uintptr(unsafe.Pointer(&src))+2, unsafe.Sizeof(byte(1))*8)
	assert.Equal(t, [11]byte{0, 0, 3, 4, 5, 6, 7, 8, 9, 10, 0}, dst)
}

func TestMemoryCopy_Array11(t *testing.T) {
	src := [13]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	dst := [13]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Do(uintptr(unsafe.Pointer(&dst))+1, uintptr(unsafe.Pointer(&src))+1, unsafe.Sizeof(byte(1))*11)
	assert.Equal(t, [13]byte{0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0}, dst)
}

func TestMemoryCopy_Array55(t *testing.T) {
	src := [64]byte{}
	dst := [64]byte{}
	expect := [64]byte{}
	for i := 0; i < 64; i++ {
		src[i] = byte(i)
		dst[i] = 0
		expect[i] = 0
	}
	for i := 0; i < 55; i++ {
		expect[7+i] = byte(7 + i)
	}
	Do(uintptr(unsafe.Pointer(&dst))+7, uintptr(unsafe.Pointer(&src))+7, unsafe.Sizeof(byte(1))*55)
	assert.Equal(t, expect, dst)
}
