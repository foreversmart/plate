package view

import (
	"fmt"
	"reflect"
)

type Field struct {
	tp     reflect.Type
	value  reflect.Value
	tag    reflect.StructTag
	isLeaf bool
	// struct, array map is a child view
	Child *View
}

func FetchField(parentTag reflect.StructTag, v reflect.Value) (f *Field, err error) {
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

	f = &Field{
		tp:     v.Type(),
		value:  v,
		tag:    parentTag,
		isLeaf: true,
	}

	switch v.Kind() {
	case reflect.Struct:
		f.isLeaf = false
		childView, err := FetchViewFromStruct(f.value)
		if err != nil {
			return nil, err
		}
		f.Child = childView
	case reflect.Map:
		f.isLeaf = false
		childView, err := FetchViewFromMap(parentTag, f.value)
		if err != nil {
			return nil, err
		}
		f.Child = childView
	case reflect.Array, reflect.Slice:
		f.isLeaf = false
		childView, err := FetchViewFromArray(parentTag, f.value)
		if err != nil {
			return nil, err
		}
		f.Child = childView

	}
	return
}
