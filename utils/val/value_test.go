package val

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestValueToString(t *testing.T) {
	var (
		vBool                  = true
		vStr                   = "string"
		vInt                   = 33
		vUint      uint        = 133
		vFloat32   float32     = 0.33
		vFloat64   float64     = 0.33333
		vInterface interface{} = 33
	)

	v := reflect.ValueOf(vBool)
	s, err := ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, "true", s)

	v = reflect.ValueOf(vStr)
	s, err = ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, "string", s)

	v = reflect.ValueOf(vInt)
	s, err = ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, "33", s)

	v = reflect.ValueOf(vUint)
	s, err = ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, "133", s)

	v = reflect.ValueOf(vFloat32)
	s, err = ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%.8f", vFloat32), s)

	v = reflect.ValueOf(vFloat64)
	s, err = ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%.8f", vFloat64), s)

	v = reflect.ValueOf(vInterface)
	s, err = ValueToString(v)
	assert.Nil(t, err)
	assert.Equal(t, "33", s)
}
