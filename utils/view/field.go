package view

import (
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

func FetchField(parentTag reflect.StructTag, full bool, v reflect.Value, tagOpt ...string) (f *Field, err error) {
	for v.Kind() == reflect.Ptr {

		if v.IsNil() {
			//return nil, fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
			f = &Field{
				tp:     v.Type(),
				value:  v,
				tag:    parentTag,
				isLeaf: true,
			}
			return
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

	if full == true {
		return
	}

	switch v.Kind() {
	case reflect.Struct:
		childView, err := FetchViewFromStruct(f.value, full, tagOpt...)
		if err != nil {
			return nil, err
		}

		f.Child = childView

	case reflect.Map:
		childView, err := FetchViewFromMap(parentTag, full, f.value)
		if err != nil {
			return nil, err
		}
		f.Child = childView
	case reflect.Array, reflect.Slice:
		childView, err := FetchViewFromArray(parentTag, full, f.value)
		if err != nil {
			return nil, err
		}
		f.Child = childView

	}

	if f.Child != nil {
		f.isLeaf = false
	}

	return
}
