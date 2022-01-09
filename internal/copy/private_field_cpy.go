package copy

import (
	"fmt"
	"reflect"
	"unsafe"
)

type typeT int

/*
DeepCopyPrivateFieldReflect copy the private field (specified by param fieldNum) from src to dst.

Both src and dst must be addressable.
*/
func (copier *Copier) DeepCopyPrivateFieldReflect(dst, src reflect.Value, fieldNum int) {
	srcField := src.Field(fieldNum)
	structField := src.Type().Field(fieldNum)
	dstPtr := unsafe.Pointer(dst.UnsafeAddr() + structField.Offset)
	fmt.Println(structField.Name)

	if isSimpleKind(srcField.Kind()) {
		srcPtr := unsafe.Pointer(src.UnsafeAddr() + structField.Offset)
		MemoryCopy(dstPtr, srcPtr, int(structField.Type.Size()))
		return
	}

	switch srcField.Kind() {
	case reflect.Ptr:
		ptr := copier.newDeepCopyOf(srcField)
		*(**typeT)(dstPtr) = *(**typeT)(unsafe.Pointer(ptr.Pointer()))

	case reflect.Slice:
		ptr := copier.newDeepCopyOf(srcField)
		*(*[]typeT)(dstPtr) = *(*[]typeT)(unsafe.Pointer(ptr.Pointer()))

	case reflect.Map:
		ptr := copier.newDeepCopyOf(srcField)
		*(*map[typeT]typeT)(dstPtr) = *(*map[typeT]typeT)(unsafe.Pointer(ptr.Pointer()))

	default:
		ptr := copier.newDeepCopyOf(srcField)
		MemoryCopy(dstPtr, unsafe.Pointer(ptr.Pointer()), int(structField.Type.Size()))
	}
}
