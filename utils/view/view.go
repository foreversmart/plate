/*
	view is a static reflect map value of a struct which means pointer value will change to real value
	default view has all the tag of a struct
	view compose with a map of field
	* field has a type and reflect value tag and so on.
*/
package view

import (
	"fmt"
	"github.com/foreversmart/plate/utils/val"
	"reflect"
	"strconv"
	"strings"
)

type View struct {
	Aspect string            // aspect is the key tag to fetch view from an object
	Fields map[string]*Field // fields means a map of filed list, key is aspect tag of filed
}

//
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

func (view *View) MergeWithNew(nv *View) {
	for k, f := range nv.Fields {
		view.Fields[k] = f
	}
}

// SetObjectValue set the object o value by the view value.
// use tag to select o aspect to set value
// use must parameter to determine ignore or occurs errors when not find the value.
func (view *View) SetObjectValue(o interface{}, must bool, tagOpt ...string) (err error) {
	v := reflect.ValueOf(o)
	v = val.SettableValue(v)
	return view.SetStructValue(v, must, false, tagOpt...)
}

func (view *View) SetStructValue(v reflect.Value, must, full bool, tagOpt ...string) (err error) {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cant fetch view from type %s", v.Type().String())
	}

	fieldNum := v.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		sf := v.Type().Field(i)
		name, isFull, isValid := fieldName(sf, full, tagOpt...)
		if !isValid {
			continue
		}

		field, ok := view.Fields[name]
		// field not match
		if !ok {
			continue
		}

		fv := v.Field(i)
		// if fv can't set jump
		if !fv.CanSet() {
			continue
		}

		if field.isLeaf || isFull {
			SetValue(fv, field.value)
			continue
		}

		fv = val.SettableValue(fv)
		switch fv.Kind() {
		case reflect.Struct:
			err = field.Child.SetStructValue(fv, must, isFull, tagOpt...)
		case reflect.Map:
			err = field.Child.SetMapValue(fv, must, isFull, tagOpt...)
		case reflect.Array, reflect.Slice:
			err = field.Child.SetArrayValue(fv, must, isFull, tagOpt...)

		}

	}

	return
}

// SetMapValue set view's value to v. v is a map or occurs type not much error
func (view *View) SetMapValue(v reflect.Value, must, full bool, tagOpt ...string) error {
	if v.Kind() != reflect.Map {
		return fmt.Errorf("cant set view type %s only support map kind", v.Type().String())
	}

	// when map is nil make a new map
	if v.IsNil() {
		ot := v.Type()
		v.Set(reflect.MakeMap(ot))
	}

	keyType := v.Type().Key()
	valueType := v.Type().Elem()

	for key, value := range view.Fields {

		keyValue, err := val.NewValueByString(key, keyType)
		if err != nil {
			return err
		}

		if value.isLeaf {
			v.SetMapIndex(keyValue, value.value)
			continue
		}

		valueValue := reflect.New(valueType)
		valueValue = val.SettableValue(valueValue)

		switch valueValue.Kind() {
		case reflect.Struct:
			err = value.Child.SetStructValue(valueValue, must, full, tagOpt...)
		case reflect.Map:
			err = value.Child.SetMapValue(valueValue, must, full, tagOpt...)
		case reflect.Array, reflect.Slice:
			err = value.Child.SetArrayValue(valueValue, must, full, tagOpt...)
		}

		v.SetMapIndex(keyValue, valueValue)

	}

	return nil
}

func (view *View) SetArrayValue(v reflect.Value, must, full bool, tagOpt ...string) (err error) {
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return fmt.Errorf("cant set view type %s only support array or slice kind", v.Kind().String())
	}

	// when array is nil make a new array
	if v.IsNil() {
		ot := v.Type()
		v.Set(reflect.MakeSlice(ot, len(view.Fields), len(view.Fields)))
	}

	for key, value := range view.Fields {
		arrayIndex, _ := strconv.ParseInt(key, 10, 64)

		if value.isLeaf {
			v.Index(int(arrayIndex)).Set(value.value)
			continue
		}

		arrayValue := val.SettableValue(v.Index(int(arrayIndex)))

		switch arrayValue.Kind() {
		case reflect.Struct:
			err = value.Child.SetStructValue(arrayValue, must, full, tagOpt...)
		case reflect.Map:
			err = value.Child.SetMapValue(arrayValue, must, full, tagOpt...)
		case reflect.Array, reflect.Slice:
			err = value.Child.SetArrayValue(arrayValue, must, full, tagOpt...)
		}
	}

	return

}
