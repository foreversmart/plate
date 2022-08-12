package val

import (
	"fmt"
	"reflect"
	"strconv"
)

func NewValueByString(s string, t reflect.Type) (reflect.Value, error) {
	return NewValueByBytes([]byte(s), t)
}

func NewValueByBytes(key []byte, t reflect.Type) (reflect.Value, error) {
	res := reflect.New(t).Elem()
	switch t.Kind() {
	case reflect.Bool:
		v, err := parseByteBool(key)
		if err != nil {
			return res, err
		}

		res.SetBool(v)
	case reflect.String:
		res.SetString(string(key))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(string(key), 10, 64)
		if err != nil {
			return res, err
		}
		res.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, err := strconv.ParseUint(string(key), 10, 64)
		if err != nil {
			return res, err
		}
		res.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(string(key), 64)
		if err != nil {
			return res, err
		}
		res.SetFloat(v)
	case reflect.Interface:
		// TODO need dynamic check the type but json never happen
		// default support as string
		res.Set(reflect.ValueOf(string(key)))
	}

	return res, nil
}

func SetValueByString(value reflect.Value, s string) error {
	if !value.CanSet() {
		return fmt.Errorf("parse %s value is not settable", s)
	}

	switch value.Kind() {
	case reflect.Bool:
		v, err := parseByteBool([]byte(s))
		if err != nil {
			return err
		}

		value.SetBool(v)
	case reflect.String:
		value.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		value.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		value.SetFloat(v)
	case reflect.Interface:
		// TODO need dynamic check the type but json never happen
		// default support as string
		value.Set(reflect.ValueOf(s))
	}

	return nil
}

func parseByteBool(b []byte) (bool, error) {
	if len(b) == 1 {
		switch b[0] {
		case '0', 'f':
			return false, nil
		case '1', 't':
			return true, nil
		}

		return false, fmt.Errorf("parse byte bool invalid %s", string(b))
	}

	if b[0] == 't' && b[1] == 'r' && b[2] == 'u' && b[3] == 'e' {
		return true, nil
	}

	if b[0] == 'f' && b[1] == 'a' && b[2] == 'l' && b[3] == 's' && b[4] == 'e' {
		return false, nil
	}

	return false, fmt.Errorf("parse byte bool invalid %s", string(b))

}

func ValueToString(v reflect.Value) (res string, err error) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		}

		return "false", nil
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		// TODO this may cause value collision
		return fmt.Sprintf("%.8f", v.Float()), nil
	case reflect.Interface:
		return fmt.Sprintf("%v", v.Interface()), nil
	default:
		return "", fmt.Errorf("cant support value kind %s", v.Kind().String())
	}

	return
}

func SettableValue(v reflect.Value) reflect.Value {
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
