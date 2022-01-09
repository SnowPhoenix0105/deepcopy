package refutils

import (
	"reflect"
	"unsafe"
)

func AddFlagFor(value *reflect.Value, flag ReflectFlag) ReflectFlag {
	ptr := (*Value)(unsafe.Pointer(value))
	ret := ptr.ReflectFlag
	ptr.ReflectFlag |= flag
	return ret
}

func DeleteFlagFor(value *reflect.Value, flag ReflectFlag) ReflectFlag {
	ptr := (*Value)(unsafe.Pointer(value))
	ret := ptr.ReflectFlag
	ptr.ReflectFlag &= ^flag
	return ret
}

func SetFlagFor(value *reflect.Value, flag ReflectFlag) {
	ptr := (*Value)(unsafe.Pointer(value))
	ptr.ReflectFlag = flag
}
