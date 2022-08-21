package view

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestView_SetObjectValue(t *testing.T) {
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

	view, err := FetchViewFromStruct(reflect.ValueOf(demo), false)
	assert.Nil(t, err)
	assert.NotNil(t, view)

	result := &TestDemo{}

	err = view.SetObjectValue(result, true)
	assert.Nil(t, err)

	assert.EqualValues(t, demo, result)

	jStr, _ := json.Marshal(result)
	fmt.Println(string(jStr))

	jStr, _ = json.Marshal(demo)
	fmt.Println(string(jStr))

	fmt.Println(result.FieldIntMap[1111])
}
