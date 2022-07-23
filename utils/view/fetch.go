package view

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	TagOption = ","
)

// FetchViewFromStruct fetch view from reflect Value v
// pass viewTags to decide which tag we use to fetch, if not passed we use the default type's name
// return view and error
func FetchViewFromStruct(v reflect.Value, viewTags ...string) (view *View, err error) {
	var (
		isSetViewTag bool
		viewTag      string
	)

	if len(viewTags) > 0 {
		viewTag = viewTags[0]
		isSetViewTag = true
	}

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
		tag := fieldType.Tag

		fieldTagName := fieldType.Name
		if isSetViewTag {
			tc := tag.Get(viewTag)
			tcl := strings.Split(tc, ",")
			fieldTagName = tcl[0]
		}

		field, fetchErr := FetchField(tag, v.Field(i))
		if fetchErr != nil {
			return nil, fetchErr
		}

		view.Fields[fieldTagName] = field
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
		field, err := FetchField(parentTag, arrValue)
		if err != nil {
			// TODO
		}

		view.Fields[fmt.Sprintf("%d", i)] = field
	}

	return
}
