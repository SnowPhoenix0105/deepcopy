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

type DeepCopyOptions struct {
	IgnoreUnexploredFields bool
}

type Copier struct {
	DeepCopyOptions
}

func (copier *Copier) DeepCopyOf(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	return copier.DeepCopyOfReflect(reflect.ValueOf(obj)).Interface()
}

func (copier *Copier) DeepCopyOfReflect(obj reflect.Value) reflect.Value {
	if !obj.IsValid() {
		return reflect.Value{}
	}
	return copier.deepCopyOf(obj)
}

func (copier *Copier) NewDeepCopyOf(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	return copier.NewDeepCopyOfReflect(reflect.ValueOf(obj)).Interface()
}

func (copier *Copier) NewDeepCopyOfReflect(obj reflect.Value) reflect.Value {
	if !obj.IsValid() {
		return reflect.Value{}
	}
	return copier.newDeepCopyOf(obj)
}

/*
newDeepCopyOf require the obj must be valid.
*/
func (copier *Copier) newDeepCopyOf(obj reflect.Value) reflect.Value {
	return copier.newDeepCopyOf2(obj, false)
}

func (copier *Copier) newDeepCopyOf2(obj reflect.Value, inDefault bool) reflect.Value {
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
		forceSet(&dst, copier.deepCopyOf2(obj, true))

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
			forceSet(&dst, copier.deepCopyOf(obj.Index(i)))
		}

	case reflect.Interface:
		if obj.IsNil() {
			break
		}
		copier.DeepCopyInterface(ptr.Elem(), obj)

	case reflect.Struct:
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
				field.Set(copier.deepCopyOf(obj.Field(i)))
			} else {
				copier.DeepCopyPrivateFieldReflect(dst, obj, i)
			}
		}
	}

	return ptr
}

/*
deepCopyOf require the obj must be valid.
*/
func (copier *Copier) deepCopyOf(obj reflect.Value) reflect.Value {
	return copier.deepCopyOf2(obj, false)
}

func (copier *Copier) deepCopyOf2(obj reflect.Value, inDefault bool) reflect.Value {
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
		return copier.newDeepCopyOf2(obj, true).Elem()

	case reflect.Ptr:
		if obj.IsNil() {
			return reflect.Zero(obj.Type())
		}
		return copier.newDeepCopyOf(obj.Elem())

	case reflect.Slice:
		if obj.IsNil() {
			return reflect.Zero(obj.Type())
		}
		length := obj.Len()
		slice := reflect.MakeSlice(obj.Type(), length, obj.Cap())
		for i := 0; i < length; i++ {
			//slice.Index(i).Set(deepCopyOf(obj.Index(i)))
			dst := slice.Index(i)
			forceSet(&dst, copier.deepCopyOf(obj.Index(i)))
		}
		return slice

	case reflect.Map:
		if obj.IsNil() {
			return reflect.Zero(obj.Type())
		}
		dict := reflect.MakeMap(obj.Type())
		iter := obj.MapRange()
		for iter.Next() {
			dict.SetMapIndex(copier.deepCopyOf(iter.Key()), copier.deepCopyOf(iter.Value()))
		}
		return dict
	}
}
