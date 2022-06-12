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
}

type Field struct {
	tp       reflect.Type
	value    reflect.Value
	tag      string
	isStruct bool
	// only support struct nested
	ChildView *View
}

func fetchViewFrom(tag string, v reflect.Value) (view *View, err error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cant fetch view from type %s", v.Type().String())
	}

	view = &View{
		Fields: make(map[string]*Field),
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		field := v.Type().Field(i)
		fieldTag := field.Tag.Get(tag)
		if fieldTag == "" || fieldTag == "-" {
			continue
		}

		vf := &Field{
			tp:    field.Type,
			value: v.Field(i),
			tag:   field.Tag.Get(tag),
		}

		if field.Type.Kind() == reflect.Struct {
			vf.isStruct = true
			childView, err := fetchViewFrom(tag, vf.value)
			if err != nil {
				return nil, err
			}
			vf.ChildView = childView
		}

	}

	return
}

func (view *View) Parse(v interface{}) {
}
