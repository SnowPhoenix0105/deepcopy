package deepcopy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopier_OfInterface(t *testing.T) {
	type Sub struct {
		integer *int
		str     [4]string
		all     interface{}
	}
	type Top struct {
		sub  Sub
		sub2 *Sub
	}

	obj := Top{
		sub: Sub{
			integer: new(int),
			str:     [4]string{"a", "", "b", "c"},
			all:     1.23,
		},
		sub2: &Sub{
			integer: new(int),
			str:     [4]string{"a", "", "b", "c"},
			all:     []int{3, 2, 1},
		},
	}
	*obj.sub.integer = 5
	*obj.sub2.integer = 3

	deepCopy := Unsafe()
	cpy, ok := deepCopy.OfInterface(obj).(Top)
	assert.True(t, ok)
	assert.Equal(t, obj, cpy)
	*obj.sub.integer = 4
	assert.Equal(t, 5, *cpy.sub.integer)
	*cpy.sub2.integer = 4
	assert.Equal(t, 3, *obj.sub2.integer)
	slice, ok := obj.sub2.all.([]int)
	assert.True(t, ok)
	slice[1] = 123
	assert.Equal(t, []int{3, 2, 1}, cpy.sub2.all)
}
