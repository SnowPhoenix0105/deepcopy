package copy

import (
	"reflect"
	"unsafe"
)

var (
	ptrSize uintptr
	//emptyInterfaceSize uintptr
	//noEmptyInterfaceSize uintptr
)

func init() {
	ptrSize = unsafe.Sizeof(uintptr(0))
	//emptyInterfaceSize = unsafe.Sizeof((interface{})(nil))
	//noEmptyInterfaceSize = unsafe.Sizeof((interface{M()})(nil))
}

/*
DeepCopyInterface deep copy an interface.
*/
func (copier *Copier) DeepCopyInterface(dst, src reflect.Value) {
	// copy the itab
	forceSet(&dst, src)
	// dst.Set(src)

	// deep-copy the word
	ptr := copier.newDeepCopyOf(src.Elem())
	//var ifaceSize uintptr
	//if src.NumMethod() == 0 {
	//	ifaceSize = emptyInterfaceSize
	//} else {
	//	ifaceSize = noEmptyInterfaceSize
	//}
	dstPtr := unsafe.Pointer(dst.UnsafeAddr() + ptrSize)
	var srcPtr unsafe.Pointer
	if src.Elem().Kind() != reflect.Ptr {
		srcPtr = unsafe.Pointer(ptr.Pointer())
	} else {
		srcPtr = unsafe.Pointer(ptr.Elem().Pointer())
	}
	*(*uintptr)(dstPtr) = uintptr(srcPtr)
}
