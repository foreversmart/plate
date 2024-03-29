package view

import (
	"fmt"
	"github.com/foreversmart/plate/utils/val"
	"reflect"
)

// FetchViewFromStruct fetch view from reflect Value v
// pass viewTags to decide which tag we use to fetch, if not passed we use the default type's name
// return view and error
func FetchViewFromStruct(v reflect.Value, full bool, tagOpt ...string) (view *View, err error) {
	// zero value cant get view
	//if v.IsZero() {
	//	return nil, nil
	//}

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			//return nil, fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
			return nil, nil
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

	if len(tagOpt) > 0 {
		view.Aspect = tagOpt[0]
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		sf := v.Type().Field(i)
		name, isFull, isValid := fieldName(sf, full, tagOpt...)
		if !isValid {
			continue
		}

		field, fetchErr := FetchField(sf.Tag, isFull, v.Field(i), tagOpt...)
		if fetchErr != nil {
			return nil, fetchErr
		}

		view.Fields[name] = field
	}

	return
}

func FetchViewFromMap(parentTag reflect.StructTag, full bool, v reflect.Value) (view *View, err error) {
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
		field, err := FetchField(parentTag, full, mapValue)
		if err != nil {
			// TODO collect the errors
		}

		keyStr, err := val.ValueToString(key)
		if err != nil {
			// TODO collect the errors
		}

		view.Fields[keyStr] = field
	}

	return
}

func FetchViewFromArray(parentTag reflect.StructTag, full bool, v reflect.Value) (view *View, err error) {
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
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
		field, err := FetchField(parentTag, full, arrValue)
		if err != nil {
			// TODO
		}

		view.Fields[fmt.Sprintf("%d", i)] = field
	}

	return
}
