package view

import (
	"github.com/valyala/fastjson"
	"reflect"
)

type View struct {
	data map[string]*fastjson.Value
}

type Field struct {
	tp      reflect.Type
	value   reflect.Value
	tag     string
	isBasic bool
	// map, array, nested struct
	Nest  map[string]*Field
	Array []*Field
	Map   map[string]*Field
}

func FetchViewFrom(tag string, v interface{}) *View {
	return nil
}

func (view *View) Parse(v interface{}) {
}
