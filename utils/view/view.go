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
	v = settableValue(v)
	return view.SetStructValue(v, must, tagOpt...)
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

		fv := v.Field(i)
		if field.isLeaf {
			SetValue(fv, field.value)
			continue
		}

		fv = settableValue(fv)
		switch fv.Kind() {
		case reflect.Struct:
			err = field.Child.SetStructValue(fv, must, tagOpt...)
		case reflect.Map:
			err = field.Child.SetMapValue(fv, must, tagOpt...)
		case reflect.Array:
			err = field.Child.SetStructValue(fv, must, tagOpt...)

		}

	}

	return
}

// SetMapValue set view's value to v. v is a map or occurs type not much error
func (view *View) SetMapValue(v reflect.Value, must bool, tagOpt ...string) error {
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
		valueValue = settableValue(valueValue)

		switch valueValue.Kind() {
		case reflect.Struct:
			err = value.Child.SetStructValue(valueValue, must, tagOpt...)
		case reflect.Map:
			err = value.Child.SetMapValue(valueValue, must, tagOpt...)
		case reflect.Array:
			err = value.Child.SetStructValue(valueValue, must, tagOpt...)
		}

		v.SetMapIndex(keyValue, valueValue)

	}

	return nil
}

func (view *View) SetArrayValue(v reflect.Value, must string, tagOpt ...string) error {

	return nil

}
