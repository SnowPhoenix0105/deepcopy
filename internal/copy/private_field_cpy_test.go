package copy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeepCopyPrivateFieldReflect(t *testing.T) {
	type Sub struct {
		str    string
		slice  []int
		parray *[2]int
	}
	type Class struct {
		integer int
		pint    *int
		array   [4]byte
		sub     Sub
	}
	var integer = 2
	obj := Class{
		integer: 1,
		pint:    &integer,
		array:   [4]byte{3, 4, 5, 6},
		sub: Sub{
			str:    "foo",
			slice:  []int{7, 8, 9},
			parray: &[2]int{10, 11},
		},
	}
	copier := Copier{
		DeepCopyOptions{
			IgnoreUnexploredFields: false,
		},
	}
	cpy := copier.DeepCopyOf(obj)
	assert.Equal(t, obj, cpy)
}
