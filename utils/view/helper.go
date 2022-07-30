package view

import (
	"reflect"
)

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
