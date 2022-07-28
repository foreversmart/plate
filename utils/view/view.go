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
	Aspect string            // aspect is the key tag to fetch view from an object
	Fields map[string]*Field // fields means a map of filed list, key is aspect tag of filed
}

func (view *View) NewView(aspect string) *View {
	newView := &View{}
	newView.Aspect = aspect
	newView.Fields = make(map[string]*Field)

	for _, f := range view.Fields {
		tc := f.tag.Get(aspect)
		tcl := strings.Split(tc, TagOption)
		fieldTagName := tcl[0]
		if !isFieldNameValid(fieldTagName) {
			continue
		}

		// leaf field
		if f.isLeaf {
			newView.Fields[fieldTagName] = f
			continue
		}

		nf := &Field{
			tp:     f.tp,
			value:  f.value,
			tag:    f.tag,
			isLeaf: false,
			Child:  f.Child.NewView(aspect),
		}
		newView.Fields[fieldTagName] = nf
	}

	return newView
}

// SetObjectValue set the object o value by the view value.
// use tag to select o aspect to set value
// use must parameter to determine ignore or occurs errors when not find the value.
func (view *View) SetObjectValue(o interface{}, must bool, tagOpt ...string) (err error) {
	v := reflect.ValueOf(o)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			ot := v.Type().Elem()
			v.Set(reflect.New(ot))
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
		sf := v.Type().Field(i)

		name, isValid := fieldName(sf, tagOpt...)
		if !isValid {
			continue
		}

		// select the right field in the view
		field, ok := view.Fields[name]
		// field not match
		if !ok {
			continue
		}

		fv := v.Field(i)
		if field.isLeaf {
			SetValue(fv, field.value)
			continue
		}

		fv = settableValue(fv)
		switch fv.Kind() {
		case reflect.Struct:
		case reflect.Map:
		case reflect.Array:

		}

		field.Child.SetStructValue(v.Field(i), must, tagOpt...)
	}

	return
}

func settableValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			ot := v.Type().Elem()
			v.Set(reflect.New(ot))
		}

		child := v.Elem()
		// break loop pointer
		if v == child {
			break
		}

		v = child
	}

	return v
}

func (view *View) SetStructValue(v reflect.Value, must bool, tagOpt ...string) (err error) {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cant fetch view from type %s", v.Type().String())
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		sf := v.Type().Field(i)
		name, isValid := fieldName(sf, tagOpt...)
		if !isValid {
			continue
		}

		field, ok := view.Fields[name]
		// field not match
		if !ok {
			continue
		}

		if field.isLeaf {
			v.Field(i).Set(field.value)
			continue
		}

		field.Child.SetStructValue(v.Field(i), must, tagOpt...)
	}

	return
}

func SetValue(obV, dataV reflect.Value) (err error) {
	for obV.Kind() == reflect.Ptr {
		if obV.IsNil() {
			// TODO optimize
			ot := obV.Type().Elem()
			obV.Set(reflect.New(ot))
			obV.Elem().Set(dataV)
			return
		}

		child := obV.Elem()
		// break loop pointer
		if obV == child {
			break
		}

		obV = child
	}

	if obV.Type() == dataV.Type() {
		obV.Set(dataV)
	}

	return nil
}

/*
A
viewMysql viewBson viewEs viewRedis viewMem viewGraph viewStruct view Raw
viewMysql.Set(b)

*/
