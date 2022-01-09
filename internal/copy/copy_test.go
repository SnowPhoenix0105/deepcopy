package copy

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStructCopy(t *testing.T) {
	type Class struct {
		privateInt int
		privateStr string
	}
	obj := Class{
		privateInt: 1,
		privateStr: "foo",
	}
	ptr := reflect.New(reflect.TypeOf(obj))
	ptr.Elem().Set(reflect.ValueOf(obj))
	assert.Equal(t, obj, ptr.Elem().Interface())
	obj.privateInt = 2
	assert.NotEqual(t, obj, ptr.Elem().Interface())
}

func TestNewDeepCopyOf_Simple(t *testing.T) {
	type MyInt int
	var integer MyInt = 123
	ptr, ok := NewDeepCopyOf(&Options{IgnoreUnexploredFields: false}, integer).(*MyInt)
	assert.True(t, ok)
	assert.Equal(t, integer, *ptr)
	*ptr = 321
	assert.Equal(t, MyInt(123), integer)
	assert.Equal(t, MyInt(321), *ptr)
}

func TestNewDeepCopyOf_Ptr(t *testing.T) {
	type MyInt int
	var integer MyInt = 123
	ptr := &integer
	cpptr, ok := NewDeepCopyOf(&Options{IgnoreUnexploredFields: false}, ptr).(**MyInt)
	assert.True(t, ok)
	assert.Equal(t, *ptr, **cpptr)
	*ptr = 234
	assert.Equal(t, MyInt(234), *ptr)
	assert.Equal(t, MyInt(123), **cpptr)
	**cpptr = 345
	assert.Equal(t, MyInt(234), *ptr)
	assert.Equal(t, MyInt(345), **cpptr)
}

func TestNewDeepCopyOf_Slice(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	ptr, ok := NewDeepCopyOf(&Options{IgnoreUnexploredFields: false}, slice).(*[]int)
	assert.True(t, ok)
	assert.Equal(t, slice, *ptr)
	(*ptr)[0] = 5
	assert.Equal(t, []int{1, 2, 3, 4}, slice)
	assert.Equal(t, []int{5, 2, 3, 4}, *ptr)
	slice = append(slice, 5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)
	assert.Equal(t, []int{5, 2, 3, 4}, *ptr)
}

func TestNewDeepCopyOf_Map(t *testing.T) {
	m := map[int]int{1: 1, 2: 2}
	ptr, ok := NewDeepCopyOf(&Options{IgnoreUnexploredFields: false}, m).(*map[int]int)
	assert.True(t, ok)
	assert.Equal(t, m, *ptr)
	(*ptr)[3] = 3
	assert.Equal(t, map[int]int{1: 1, 2: 2}, m)
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3}, *ptr)
}

func TestNewDeepCopyOf_MapPtr(t *testing.T) {
	array := [...]int{0, 1, 2, 3, 4}
	m := map[int]*int{1: &array[1], 2: &array[2]}
	ptr, ok := NewDeepCopyOf(&Options{IgnoreUnexploredFields: false}, m).(*map[int]*int)
	assert.True(t, ok)
	assert.False(t, m[1] == (*ptr)[1])
	assert.Equal(t, *(m[1]), *((*ptr)[1]))
}
