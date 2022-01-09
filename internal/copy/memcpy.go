package copy

import "unsafe"

const MaxMemoryCopyLength int = 1 << 20

func MemoryCopy(dst, src unsafe.Pointer, byteCount int) {
	copy((*(*[MaxMemoryCopyLength]byte)(dst))[:byteCount], (*(*[MaxMemoryCopyLength]byte)(src))[:byteCount])
}
