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
