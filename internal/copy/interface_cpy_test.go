package copy

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unsafe"
)

func TestInterfaceSize(t *testing.T) {
	emptySize := unsafe.Sizeof((interface{})(nil))
	notemptySize := unsafe.Sizeof((interface{ M() })(nil))
	t.Log(emptySize)
	assert.Equal(t, emptySize, notemptySize)
	assert.Equal(t, ptrSize*2, emptySize)
	assert.Equal(t, ptrSize*2, notemptySize)
}

func TestEmptyInterfaceCopy(t *testing.T) {
	integers := [...]int{0, 1, 2, 3, 4}
	var iface interface{} = integers[1]
	dstPtr := unsafe.Pointer(reflect.ValueOf(&iface).Pointer() + ptrSize)
	srcPtr := unsafe.Pointer(reflect.ValueOf(&integers[2]).Pointer())
	*(*uintptr)(dstPtr) = (uintptr)(srcPtr)

	assert.Equal(t, 2, iface)
}

type testInteger int

func (i testInteger) M() int {
	return int(i)
}

type testInteger2 int

func (i *testInteger2) M() int {
	return int(*i)
}

func TestNoEmptyInterfaceCopy(t *testing.T) {
	integers := [...]testInteger{0, 1, 2, 3, 4}
	var iface interface{ M() int } = integers[1]
	dstPtr := unsafe.Pointer(reflect.ValueOf(&iface).Pointer() + ptrSize)
	srcPtr := unsafe.Pointer(reflect.ValueOf(&integers[2]).Pointer())
	*(*uintptr)(dstPtr) = (uintptr)(srcPtr)

	assert.Equal(t, testInteger(2), iface)
}

func TestDeepCopyInterface(t *testing.T) {
	type iface interface{ M() int }
	type Class struct {
		Interface iface
	}
	copier := Copier{}
	obj := Class{Interface: testInteger(1)}
	cpy := copier.DeepCopyOf(obj).(Class)
	assert.Equal(t, obj, cpy)
	assert.Equal(t, 1, cpy.Interface.M())
	obj.Interface = testInteger(2)
	assert.Equal(t, 1, cpy.Interface.M())
}

func TestDeepCopyInterface2(t *testing.T) {
	type iface interface{ M() int }
	type Class struct {
		Interface iface
	}
	copier := Copier{}
	array := [...]testInteger2{0, 1, 2, 3}
	obj := &Class{Interface: &array[1]}
	cpy := copier.DeepCopyOf(obj).(*Class)
	assert.Equal(t, obj, cpy)
	assert.Equal(t, 1, cpy.Interface.M())
	obj.Interface = &array[2]
	assert.Equal(t, 1, cpy.Interface.M())
}
