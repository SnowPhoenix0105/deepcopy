package memcpy

import (
	"unsafe"
)

func Do(dst, src, byteCount uintptr) {
	switch dst % 8 {
	case 1:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 2:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 3:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 4:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 5:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 6:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 7:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		byteCount--
		fallthrough
	case 0:
	}
	for byteCount >= 8 {
		*(*uint64)(unsafe.Pointer(dst)) = *(*uint64)(unsafe.Pointer(src))
		dst += 8
		src += 8
		byteCount -= 8
	}
	switch byteCount % 8 {
	case 7:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 6:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 5:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 4:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 3:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 2:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 1:
		*(*byte)(unsafe.Pointer(dst)) = *(*byte)(unsafe.Pointer(src))
		dst++
		src++
		fallthrough
	case 0:
	}
}
