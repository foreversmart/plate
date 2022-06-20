package view

import (
	"fmt"
	"reflect"
)

type View struct {
	Fields map[string]*Field
}

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

type Field struct {
	tp     reflect.Type
	value  reflect.Value
	tag    reflect.StructTag
	isLeaf bool
	// struct, array map is a child view
	Child *View
}

type A struct {
	a *A
}

func fetchViewFromStruct(v reflect.Value) (view *View, err error) {
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
		field, err := fetchField(v.Type().Field(i).Tag, v.Field(i))
	}

	return
}

func fetchField(parentTag reflect.StructTag, v reflect.Value) (f *Field, err error) {
	vf := &Field{
		tp:     v.Type(),
		value:  v,
		tag:    parentTag,
		isLeaf: true,
	}

	switch v.Kind() {
	case reflect.Struct:
		vf.isLeaf = false
		childView, err := fetchViewFromStruct(vf.value)
		if err != nil {
			return nil, err
		}
		vf.Child = childView
	case reflect.Map:
		vf.isLeaf = false
		childView, err := fetchViewFromMap(parentTag, vf.value)
		if err != nil {
			return nil, err
		}
		vf.Child = childView
	case reflect.Array:
		vf.isLeaf = false
		childView, err := fetchViewFromMap(parentTag, vf.value)
		if err != nil {
			return nil, err
		}
		vf.Child = childView

	}
	return
}

func fetchViewFromMap(parentTag reflect.StructTag, v reflect.Value) (view *View, err error) {
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
		field, err := fetchField(parentTag, mapValue)
		if err != nil {
			// TODO
		}

		view.Fields[key.String()] = field
	}

	return
}

func fetchViewFromArray(parentTag reflect.StructTag, v reflect.Value) (view *View, err error) {
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
		field, err := fetchField(parentTag, arrValue)
		if err != nil {
			// TODO
		}

		view.Fields[fmt.Sprintf("%d", i)] = field
	}

	return
}

func (view *View) Parse(v interface{}) {
}
