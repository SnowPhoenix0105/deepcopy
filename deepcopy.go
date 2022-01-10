package deepcopy

import (
	copy2 "github.com/SnowPhoenix0105/deepcopy/internal/copy"
	"reflect"
)

type Options = copy2.Options
type UnsafeOptions = copy2.UnsafeOptions

type copier struct {
	options *Options
}

type Copier interface {
	OfInterface(obj interface{}) interface{}
	OfReflect(obj reflect.Value) reflect.Value
	AddressableOfReflect(obj reflect.Value) reflect.Value
	// Of[T any](obj T) T
}

var defaultCopier = copier{
	options: NewDefaultOption(),
}

var fullCopier = copier{
	options: &copy2.Options{
		CreateNewChan: false,
		IgnoreFunc:    false,
		Unsafe: copy2.UnsafeOptions{
			DeepCopyInterface:        true,
			DeepCopyUnexportedFields: true,
		},
	},
}

func NewDefaultOption() *Options {
	return &copy2.Options{
		CreateNewChan: false,
		IgnoreFunc:    false,
		Unsafe: copy2.UnsafeOptions{
			DeepCopyUnexportedFields: false,
			DeepCopyInterface:        false,
		},
	}
}

func WithOptions(options *Options) Copier {
	return &copier{
		options: options,
	}
}

func SetDefaultOptions(options *Options) {
	defaultCopier = copier{
		options: options,
	}
}

func Default() Copier {
	return &defaultCopier
}

func Unsafe() Copier {
	return &fullCopier
}

func (copier *copier) OfInterface(obj interface{}) interface{} {
	return copy2.DeepCopyOf(copier.options, obj)
}

func (copier *copier) OfReflect(obj reflect.Value) reflect.Value {
	return copy2.DeepCopyOfReflect(copier.options, obj)
}

func (copier *copier) AddressableOfReflect(obj reflect.Value) reflect.Value {
	return copy2.NewDeepCopyOfReflect(copier.options, obj).Elem()
}

func OfInterface(obj interface{}) interface{} {
	return defaultCopier.OfInterface(obj)
}

func OfReflect(obj reflect.Value) reflect.Value {
	return defaultCopier.OfReflect(obj)
}

func AddressableOfReflect(obj reflect.Value) reflect.Value {
	return defaultCopier.AddressableOfReflect(obj)
}
