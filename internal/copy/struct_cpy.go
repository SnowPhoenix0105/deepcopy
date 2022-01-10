package copy

import (
	"github.com/SnowPhoenix0105/deepcopy/pkg/memcpy"
	"reflect"
	"unsafe"
)

type typeT int

func deepCopyExportedField(options *Options, ptr reflect.Value, obj reflect.Value) {
	numField := obj.NumField()
	dst := ptr.Elem()
	for i := 0; i < numField; i++ {
		field := dst.Field(i)
		if field.CanSet() {
			field.Set(deepCopyOf(options, obj.Field(i)))
		}
	}
}

func deepCopyAllField(options *Options, ptr reflect.Value, obj reflect.Value) {
	numField := obj.NumField()
	dst := ptr.Elem()
	// dst.Set(obj)
	forceSet(&dst, obj)
	for i := 0; i < numField; i++ {
		field := dst.Field(i)
		if isSimpleKind(field.Kind()) {
			// it has already been copied by dst.Set(obj)
			continue
		}
		if field.CanSet() {
			field.Set(deepCopyOf(options, obj.Field(i)))
		} else {
			DeepCopyPrivateFieldReflect(options, dst, obj, i)
		}
	}
}

/*
DeepCopyPrivateFieldReflect copy the private field (specified by param fieldNum) from src to dst.

Both src and dst must be addressable.
*/
func DeepCopyPrivateFieldReflect(options *Options, dst, src reflect.Value, fieldNum int) {
	srcField := src.Field(fieldNum)
	structField := src.Type().Field(fieldNum)
	dstAddr := dst.UnsafeAddr() + structField.Offset
	// fmt.Println(structField.Name)

	if isSimpleKind(srcField.Kind()) {
		srcPtr := src.UnsafeAddr() + structField.Offset
		memcpy.Do(dstAddr, srcPtr, structField.Type.Size())
		return
	}

	switch srcField.Kind() {
	case reflect.Ptr:
		ptr := newDeepCopyOf(options, srcField)
		*(**typeT)(unsafe.Pointer(dstAddr)) = *(**typeT)(unsafe.Pointer(ptr.Pointer()))

	case reflect.Slice:
		ptr := newDeepCopyOf(options, srcField)
		*(*[]typeT)(unsafe.Pointer(dstAddr)) = *(*[]typeT)(unsafe.Pointer(ptr.Pointer()))

	case reflect.Map:
		ptr := newDeepCopyOf(options, srcField)
		*(*map[typeT]typeT)(unsafe.Pointer(dstAddr)) = *(*map[typeT]typeT)(unsafe.Pointer(ptr.Pointer()))

	default:
		ptr := newDeepCopyOf(options, srcField)
		memcpy.Do(dstAddr, ptr.Pointer(), structField.Type.Size())
	}
}
