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
	"strings"
)

type View struct {
	Fields map[string]*Field
}

// SetObjectValue set the object o value by the view value.
// use tag to select o aspect to set value
// use must parameter to determine ignore or occurs errors when not find the value.
func (view *View) SetObjectValue(o interface{}, tag string, must bool) (err error) {
	v := reflect.ValueOf(o)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
		}

		child := v.Elem()
		// break loop pointer
		if v == child {
			break
		}

		v = child
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cant fetch view from type %s", v.Type().String())
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		fieldType := v.Type().Field(i)
		tagStr := fieldType.Tag.Get(tag)
		items := strings.Split(tagStr, TagOption)
		if len(items) == 0 || items[0] == "-" {
			continue
		}

		field := view.Fields[items[0]]
		if field.isLeaf {
			v.Field(i).Set(field.value)
			continue
		}

		field.Child.SetRawValue(v.Field(i), tag, must)
	}

	return
}

func (view *View) SetRawValue(v reflect.Value, tag string, must bool) (err error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
		}

		child := v.Elem()
		// break loop pointer
		if v == child {
			break
		}

		v = child
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cant fetch view from type %s", v.Type().String())
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		fieldType := v.Type().Field(i)
		tagStr := fieldType.Tag.Get(tag)
		items := strings.Split(tagStr, TagOption)
		if len(items) == 0 || items[0] == "-" {
			continue
		}

		field := view.Fields[items[0]]
		if field.isLeaf {
			v.Field(i).Set(field.value)
			continue
		}

		field.Child.SetRawValue(v.Field(i), tag, must)
	}

	return
}
