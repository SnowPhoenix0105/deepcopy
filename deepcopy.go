package deepcopy

import (
	copy2 "github.com/SnowPhoenix0105/deepcopy/internal/copy"
	"reflect"
)

type Options = copy2.Options

type copier struct {
	options *Options
}

type Copier interface {
	OfInterface(obj interface{}) interface{}
	OfReflect(obj reflect.Value) reflect.Value
	// Of[T any](obj T) T
}

var defaultCopier = copier{
	options: &Options{
		IgnoreUnexploredFields: false,
	},
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

func (copier *copier) OfInterface(obj interface{}) interface{} {
	return copy2.DeepCopyOf(copier.options, obj)
}

func (copier *copier) OfReflect(obj reflect.Value) reflect.Value {
	return copy2.DeepCopyOfReflect(copier.options, obj)
}

func OfInterface(obj interface{}) interface{} {
	return defaultCopier.OfInterface(obj)
}

func OfReflect(obj reflect.Value) reflect.Value {
	return defaultCopier.OfReflect(obj)
}
