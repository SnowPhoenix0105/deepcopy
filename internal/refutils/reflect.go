package refutils

import "unsafe"

type ReflectFlag uintptr

type ReflectTflag uint8
type ReflectNameOff int32 // offset to a name
type ReflectTypeOff int32 // offset to an *ReflectRtype

const (
	ReflectFlagKindWidth               = 5 // there are 27 kinds
	ReflectFlagKindMask    ReflectFlag = 1<<ReflectFlagKindWidth - 1
	ReflectFlagStickyRO    ReflectFlag = 1 << 5
	ReflectFlagEmbedRO     ReflectFlag = 1 << 6
	ReflectFlagIndir       ReflectFlag = 1 << 7
	ReflectFlagAddr        ReflectFlag = 1 << 8
	ReflectFlagMethod      ReflectFlag = 1 << 9
	ReflectFlagMethodShift             = 10
	ReflectFlagRO          ReflectFlag = ReflectFlagStickyRO | ReflectFlagEmbedRO
)

// Value is the reflection interface to a Go value.
//
// Not all methods apply to all kinds of values. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of value before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run time panic.
//
// The zero Value represents no value.
// Its IsValid method returns false, its Kind method returns Invalid,
// its String method returns "<invalid Value>", and all other methods panic.
// Most functions and methods never return an invalid value.
// If one does, its documentation states the conditions explicitly.
//
// A Value can be used concurrently by multiple goroutines provided that
// the underlying Go value can be used concurrently for the equivalent
// direct operations.
//
// To compare two Values, compare the results of the Interface method.
// Using == on two Values does not compare the underlying values
// they represent.
type Value struct {
	// typ holds the type of the value represented by a Value.
	Typ *ReflectRtype

	// Pointer-valued data or, if flagIndir is set, pointer to data.
	// Valid when either flagIndir is set or typ.pointers() is true.
	Ptr unsafe.Pointer

	// flag holds metadata about the value.
	// The lowest bits are flag bits:
	//	- flagStickyRO: obtained via unexported not embedded field, so read-only
	//	- flagEmbedRO: obtained via unexported embedded field, so read-only
	//	- flagIndir: val holds a pointer to the data
	//	- flagAddr: v.CanAddr is true (implies flagIndir)
	//	- flagMethod: v is a method value.
	// The next five bits give the Kind of the value.
	// This repeats typ.Kind() except for method values.
	// The remaining 23+ bits give a method number for method values.
	// If flag.kind() != Func, code can assume that flagMethod is unset.
	// If ifaceIndir(typ), code can assume that flagIndir is set.
	ReflectFlag

	// A method value represents a curried method invocation
	// like r.Read for some receiver r. The typ+val+flag bits describe
	// the receiver r, but the flag's Kind bits say Func (methods are
	// functions), and the top bits of the flag give the method number
	// in r's type's method table.
}

// ReflectRtype is the common implementation of most values.
// It is embedded in other struct types.
//
// ReflectRtype must be kept in sync with ../runtime/type.go:/^type._type.
type ReflectRtype struct {
	Size       uintptr
	Ptrdata    uintptr      // number of bytes in the type that can contain pointers
	Hash       uint32       // hash of type; avoids computation in hash tables
	Tflag      ReflectTflag // extra type information flags
	Align      uint8        // alignment of variable with this type
	FieldAlign uint8        // alignment of struct field with this type
	Kind       uint8        // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal     func(unsafe.Pointer, unsafe.Pointer) bool
	Gcdata    *byte          // garbage collection data
	Str       ReflectNameOff // string form
	PtrToThis ReflectTypeOff // type for pointer to this type, may be zero
}
