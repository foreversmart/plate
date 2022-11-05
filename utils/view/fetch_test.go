package view

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type demo1 struct {
	a    string "a"
	nest struct {
		p int `json:"p"`
	} `json:"nest"`
}

type demo2 struct {
	a    string "a"
	p    string `json:"p"`
	nest struct {
		p int `json:"p"`
	} `json:"nest"`
	arr []struct {
		p int `json:"p"`
	} `json:"arr"`
	maps map[string]struct {
	}
}

type A struct {
	a *A
}

type TestDemo struct {
	FieldString      string              `json:"field_string" plate:"field_string"`
	FieldInt         int                 `json:"field_int" plate:"field_int"`
	FieldIntPointer  *int                `json:"field_int_pointer" plate:"field_int_pointer"`
	FieldFloat       float64             `json:"field_float" plate:"field_float"`
	FieldBool        bool                `json:"field_bool" plate:"field_bool"`
	FieldNest        TestNest            `json:"field_nest" plate:"field_nest"`
	FieldNestPointer *TestNest           `json:"field_nest_pointer" plate:"field_nest_pointer"`
	FieldMap         map[string]string   `json:"field_map" plate:"field_map"`
	FieldIntMap      map[int]string      `json:"field_int_map" plate:"field_int_map"`
	FieldStructMap   map[string]TestNest `json:"field_struct_map" plate:"field_struct_map"`
	FieldArray       []string            `json:"field_array" plate:"field_array"`
	FieldIntArray    []int               `json:"field_int_array" plate:"field_int_array"`
	FieldTime        time.Time           `json:"field_time" plate:"field_time,mid:full"`
}

type TestV struct {
	Field int `json:"field_int"`
}

type TestNest struct {
	NestString string
}

func TestFetchViewFromStruct(t *testing.T) {
	var i = 10
	demo := &TestDemo{
		FieldString:     "string",
		FieldInt:        20,
		FieldIntPointer: &i,
		FieldFloat:      32.222,
		FieldBool:       true,
		FieldNest: TestNest{
			NestString: "nest string",
		},
		FieldNestPointer: &TestNest{
			NestString: "nest string 1",
		},
		FieldMap:       map[string]string{"1": "value 1", "2": "value 2"},
		FieldIntMap:    map[int]string{1: "value 1", 2: "value 2"},
		FieldStructMap: map[string]TestNest{"first": {"1"}, "second": {"2"}},
		FieldArray:     []string{"1", "2", "3"},
		FieldIntArray:  []int{1, 2, 3},
		FieldTime:      time.Now(),
	}

	view, err := FetchViewFromStruct(reflect.ValueOf(demo), false, "plate")
	assert.Nil(t, err)
	for k, v := range view.Fields {
		switch k {
		case "field_string":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldString, v.value.String())
		case "field_int":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldInt, v.value.Int())
		case "field_int_pointer":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, *demo.FieldIntPointer, v.value.Int())
		case "field_float":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldFloat, v.value.Float())
		case "field_nest":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldNest, v.value.Interface())
		case "field_nest_pointer":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, *demo.FieldNestPointer, v.value.Interface())
		case "field_map":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldMap, v.value.Interface())
		case "field_int_map":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldIntMap, v.value.Interface())
		case "field_struct_map":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldStructMap, v.value.Interface())
		case "field_array":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldArray, v.value.Interface())
		case "field_int_array":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldIntArray, v.value.Interface())
		case "field_time":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldTime, v.value.Interface())
		}
		fmt.Println(k, v, v.isLeaf)
	}
}
