package tag

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
	"reflect"
	"testing"
)

type jsonStruct struct {
	Num              int                    `json:"num" plate:"num"`
	NumUint          uint                   `json:"num_uint" plate:"num_uint"`
	Bool             bool                   `json:"bool" plate:"bool"`
	String           string                 `json:"string" plate:"string"`
	Float64          float64                `json:"float64" plate:"float64"`
	Float32          float32                `json:"float32" plate:"float32"`
	Interface        interface{}            `json:"interface" plate:"interface"`
	InterfaceFloat   interface{}            `json:"interface_float" plate:"interface_float"`
	InterfaceBool1   interface{}            `json:"interface_bool_1" plate:"interface_bool_1"`
	InterfaceBool2   interface{}            `json:"interface_bool_2" plate:"interface_bool_2"`
	InterfaceArray   interface{}            `json:"interface_array" plate:"interface_array"`
	InterfaceArrayOb interface{}            `json:"interface_array_ob" plate:"interface_array_ob"`
	InterfaceOb      interface{}            `json:"interface_ob" plate:"interface_ob"`
	Map              map[string]int         `json:"map" plate:"map"`
	MapOb            map[string]*jsonStruct `json:"map_ob" plate:"map_ob"`
	Array            []string               `json:"array" plate:"array"`
	ArrayOb          []*jsonStruct          `json:"array_ob" plate:"array_ob"`
	Child            *jsonStruct            `json:"child" plate:"child"`
}

var testJsonStr = `{
		"num": 10,
		"num_uint": 101,
		"bool": true,
		"string": "hello",
		"float64": 3.141586,
		"float32": 3.141586211,
		"interface": 333,
		"interface_float": 333.333,
		"interface_bool_1": true,
		"interface_bool_2": false,
		"interface_array": [1, 2, 3],
		"interface_array_ob": [{"one": 1, "two": "string"}],
		"interface_ob": {"one": 1, "two": "string"},
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
		"array": ["first", "second", "third"],
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

	Parse(reValue.Elem(), jsonValue, nil, nil)
	newI := reValue.Elem().Interface()

	fmt.Println(newI)
	resultStr, err := json.Marshal(newI)
	assert.Nil(t, err)

	var expectObject jsonStruct
	json.Unmarshal([]byte(testJsonStr), &expectObject)

	assert.EqualValues(t, expectObject, newI)

	fmt.Println(string(resultStr))

}

func TestReflect(t *testing.T) {
	var b jsonStruct
	// -----
	var a interface{}

	var c interface{}
	a = &c

	var x float64 = 3.4
	v := reflect.ValueOf(&x).Elem()
	v.SetFloat(7.1)
	fmt.Println(x, "x")

	vb := reflect.ValueOf(b).Field(6)
	va := reflect.ValueOf(&a).Elem()
	vd := reflect.ValueOf(b)
	fmt.Println("settability of v:", v.CanSet(), vb.CanSet(), va.CanSet(), vd.CanSet())
	fmt.Println(vb, va, va.Type())
	va.Set(reflect.ValueOf(1))
	fmt.Println(a)

}
