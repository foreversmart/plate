package tagger

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
	"reflect"
	"testing"
)

type jsonStruct struct {
	Num     int                    `json:"num" plate:"num"`
	NumUint uint                   `json:"num_uint" plate:"num_uint"`
	Bool    bool                   `json:"bool" plate:"bool"`
	String  string                 `json:"string" plate:"string"`
	Float64 float64                `json:"float64" plate:"float64"`
	Float32 float32                `json:"float32" plate:"float32"`
	Map     map[string]int         `json:"map" plate:"map"`
	MapOb   map[string]*jsonStruct `json:"map_ob" plate:"map_ob"`
	Array   []string               `json:"array" plate:"array"`
	ArrayOb []*jsonStruct          `json:"array_ob" plate:"array_ob"`
	Child   *jsonStruct            `json:"child" plate:"child"`
}

var testJsonStr = `{
		"num": 10,
		"num_uint": 101,
		"bool": true,
		"string": "hello",
		"float64": 3.141586,
		"float32": 3.141586211,
		"map": {
			"first": 20,
			"second": 30,
			"third": 50,
			"fourth": 1111111111111110
		},
		"map_ob": {
			"first": {
				"num": 10,
				"num_uint": 101,
				"bool": true,
				"string": "hello",
				"float64": 3.141586,
				"float32": 3.141586211
			}
		},
		"child": {
			"num": 10,
			"num_uint": 101,
			"bool": true,
			"string": "hello",
			"float64": 3.141586,
			"float32": 3.141586211
		}
}`

var expectValue = &jsonStruct{
	Num:    10,
	String: "hello",
}

func TestParse(t *testing.T) {
	var i interface{}
	i = jsonStruct{}
	reValue := reflect.New(reflect.TypeOf(i))
	jsonValue, err := fastjson.Parse(testJsonStr)
	if err != nil {
		panic(err)
	}

	parse(reValue.Elem(), jsonValue, nil)
	newI := reValue.Elem().Interface()

	fmt.Println(newI)
	resultStr, err := json.Marshal(newI)
	assert.Nil(t, err)

	var expectObject jsonStruct
	json.Unmarshal([]byte(testJsonStr), &expectObject)

	assert.EqualValues(t, expectObject, newI)

	fmt.Println(string(resultStr))

}
