package view

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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
	FieldString      string
	FieldInt         int `json:"field_int"`
	FieldIntPointer  *int
	FieldFloat       float64
	FieldBool        bool
	FieldNest        TestNest
	FieldNestPointer *TestNest
	FieldMap         map[string]string
	FieldIntMap      map[int]string
	FieldStructMap   map[string]TestNest
	FieldArray       []string
	FieldIntArray    []int
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
	}

	view, err := FetchViewFromStruct(reflect.ValueOf(demo))
	assert.Nil(t, err)
	for k, v := range view.Fields {
		switch k {
		case "FieldString":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldString, v.value.String())
		case "FieldInt":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldInt, v.value.Int())
		case "FieldIntPointer":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, *demo.FieldIntPointer, v.value.Int())
		case "FieldFloat":
			assert.True(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldFloat, v.value.Float())
		case "FieldNest":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldNest, v.value.Interface())
		case "FieldNestPointer":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, *demo.FieldNestPointer, v.value.Interface())
		case "FieldMap":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldMap, v.value.Interface())
		case "FieldIntMap":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldIntMap, v.value.Interface())
		case "FieldStructMap":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldStructMap, v.value.Interface())
		case "FieldArray":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldArray, v.value.Interface())
		case "FieldIntArray":
			assert.False(t, v.isLeaf)
			assert.EqualValues(t, demo.FieldIntArray, v.value.Interface())
		}
		fmt.Println(k, v)
	}
}
