/*
	view is a static reflect map value of a struct which means pointer value will change to real value
	default view has all the tag of a struct
	view compose with a map of field
	* field has a type and reflect value tag and so on.
*/
package view

import (
	"fmt"
	"reflect"
)

type View struct {
	Fields map[string]*Field
}

func FetchViewFromStruct(v reflect.Value) (view *View, err error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
		}

		child := v.Elem()
		// break loop pointer
		if v == child {
			break
		}

		v = child
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cant fetch view from type %s", v.Type().String())
	}

	view = &View{
		Fields: make(map[string]*Field),
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		fieldType := v.Type().Field(i)
		field, fetchErr := FetchField(fieldType.Tag, v.Field(i))
		if fetchErr != nil {
			return nil, fetchErr
		}

		view.Fields[fieldType.Name] = field
	}

	return
}

func FetchViewFromMap(parentTag reflect.StructTag, v reflect.Value) (view *View, err error) {
	if v.Kind() != reflect.Map {
		return nil, fmt.Errorf("fetchViewFromMap from %s", v.Kind().String())
	}

	view = &View{}
	keys := v.MapKeys()
	if len(keys) > 0 {
		view.Fields = make(map[string]*Field)
	}

	for _, key := range keys {
		mapValue := v.MapIndex(key)
		field, err := FetchField(parentTag, mapValue)
		if err != nil {
			// TODO
		}

		view.Fields[key.String()] = field
	}

	return
}

func FetchViewFromArray(parentTag reflect.StructTag, v reflect.Value) (view *View, err error) {
	if v.Kind() != reflect.Array {
		return nil, fmt.Errorf("fetchViewFromArray from %s", v.Kind().String())
	}

	view = &View{}
	l := v.Len()
	if l == 0 {
		return nil, nil
	}

	view.Fields = make(map[string]*Field)

	for i := 0; i < l; i++ {
		arrValue := v.Index(i)
		field, err := FetchField(parentTag, arrValue)
		if err != nil {
			// TODO
		}

		view.Fields[fmt.Sprintf("%d", i)] = field
	}

	return
}

func (view *View) Parse(v interface{}) {
}
