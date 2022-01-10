package copy

import (
	"errors"
	"fmt"
	"github.com/SnowPhoenix0105/deepcopy/internal/refutils"
	"reflect"
	"unsafe"
)

func forceSet(dst *reflect.Value, src reflect.Value) {
	old := refutils.DeleteFlagFor(&src, refutils.ReflectFlagRO)
	dst.Set(src)
	// refutils.SetFlagFor(&src, old)
	_ = old
}

func panicUnsupported(typ reflect.Type) {
	panic(errors.New(fmt.Sprintf("unsuport type: %s:%s", typ.PkgPath(), typ.Name())))
}

func isSimpleKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool,
		reflect.Complex64, reflect.Complex128,
		reflect.String,
		reflect.Uintptr, reflect.UnsafePointer:

		return true
	default:
		return false
	}
}

func simpleKindAssign(dst, src reflect.Value) {
	switch dst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dst.SetInt(src.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dst.SetUint(src.Uint())

	case reflect.Float32, reflect.Float64:
		dst.SetFloat(src.Float())

	case reflect.Bool:
		dst.SetBool(src.Bool())

	case reflect.Complex64, reflect.Complex128:
		dst.SetComplex(src.Complex())

	case reflect.String:
		dst.SetString(src.String())

	case reflect.Uintptr:
		dst.SetUint(src.Uint())

	case reflect.UnsafePointer:
		dst.SetPointer(unsafe.Pointer(src.Pointer()))
	}
}

type Options struct {
	CreateNewChan bool
	IgnoreFunc    bool
	Unsafe        UnsafeOptions
}

type UnsafeOptions struct {
	DeepCopyUnexportedFields bool
	DeepCopyInterface        bool
}

func DeepCopyOf(options *Options, obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	return DeepCopyOfReflect(options, reflect.ValueOf(obj)).Interface()
}

func DeepCopyOfReflect(options *Options, obj reflect.Value) reflect.Value {
	if !obj.IsValid() {
		return reflect.Value{}
	}
	return deepCopyOf(options, obj)
}

func NewDeepCopyOf(options *Options, obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	return NewDeepCopyOfReflect(options, reflect.ValueOf(obj)).Interface()
}

func NewDeepCopyOfReflect(options *Options, obj reflect.Value) reflect.Value {
	if !obj.IsValid() {
		return reflect.Value{}
	}
	return newDeepCopyOf(options, obj)
}

/*
newDeepCopyOf require the obj must be valid.
*/
func newDeepCopyOf(options *Options, obj reflect.Value) reflect.Value {
	return newDeepCopyOf2(options, obj, false)
}

func newDeepCopyOf2(options *Options, obj reflect.Value, inDefault bool) reflect.Value {
	ptr := reflect.New(obj.Type())

	if isSimpleKind(obj.Kind()) {
		simpleKindAssign(ptr.Elem(), obj)
		return ptr
	}

	switch obj.Kind() {
	default:
		if inDefault {
			panicUnsupported(obj.Type())
		}
		dst := ptr.Elem()
		forceSet(&dst, deepCopyOf2(options, obj, true))

	case reflect.Array:
		if isSimpleKind(obj.Type().Elem().Kind()) {
			//ptr.Elem().Set(obj)
			dst := ptr.Elem()
			forceSet(&dst, obj)
			break
		}
		length := obj.Len()
		array := ptr.Elem()
		for i := 0; i < length; i++ {
			//array.Index(i).Set(deepCopyOf(obj.Index(i)))
			dst := array.Index(i)
			forceSet(&dst, deepCopyOf(options, obj.Index(i)))
		}

	case reflect.Interface:
		if obj.IsNil() {
			break
		}
		if options.Unsafe.DeepCopyInterface {
			DeepCopyInterface(options, ptr.Elem(), obj)
		} else {
			ptr.Elem().Set(obj)
		}

	case reflect.Struct:
		if options.Unsafe.DeepCopyUnexportedFields {
			deepCopyAllField(options, ptr, obj)
		} else {
			deepCopyExportedField(options, ptr, obj)
		}
	}

	return ptr
}

/*
deepCopyOf require the obj must be valid.
*/
func deepCopyOf(options *Options, obj reflect.Value) reflect.Value {
	return deepCopyOf2(options, obj, false)
}

func deepCopyOf2(options *Options, obj reflect.Value, inDefault bool) reflect.Value {
	if isSimpleKind(obj.Kind()) {
		return obj
	}
	switch obj.Kind() {
	default:
		if inDefault {
			if inDefault {
				panicUnsupported(obj.Type())
			}
		}
		return newDeepCopyOf2(options, obj, true).Elem()

	case reflect.Ptr:
		if obj.IsNil() {
			return reflect.Zero(obj.Type())
		}
		return newDeepCopyOf(options, obj.Elem())

	case reflect.Slice:
		if obj.IsNil() {
			return reflect.Zero(obj.Type())
		}
		length := obj.Len()
		slice := reflect.MakeSlice(obj.Type(), length, obj.Cap())
		for i := 0; i < length; i++ {
			//slice.Index(i).Set(deepCopyOf(obj.Index(i)))
			dst := slice.Index(i)
			forceSet(&dst, deepCopyOf(options, obj.Index(i)))
		}
		return slice

	case reflect.Map:
		if obj.IsNil() {
			return reflect.Zero(obj.Type())
		}
		dict := reflect.MakeMap(obj.Type())
		iter := obj.MapRange()
		for iter.Next() {
			dict.SetMapIndex(deepCopyOf(options, iter.Key()), deepCopyOf(options, iter.Value()))
		}
		return dict

	case reflect.Chan:
		if options.CreateNewChan {
			return reflect.MakeChan(obj.Type(), obj.Len())
		} else {
			return obj
		}

	case reflect.Func:
		if options.IgnoreFunc {
			return reflect.Zero(obj.Type())
		} else {
			return obj
		}
	}
}
