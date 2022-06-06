package tagger

import (
	"fmt"
	"github.com/valyala/fastjson"
	"reflect"
	"strconv"
	"strings"
)

//func Parse(request *http.Request, v reflect.Value, prefix string) error {
//	body, err := ioutil.ReadAll(request.Body)
//	defer request.Body.Close()
//	if err != nil {
//		return err
//	}
//
//	jsonValue, err := fastjson.Parse(string(body))
//	if err != nil {
//		return err
//	}
//
//}

const (
	TagNameFetch = "plate"
	TagNameCheck = "check"
	TagOption    = ","
)

/*
	type A struct {
		B string `plate:"b,body" check:"int>10"`
		B string `plate:"b,header"`
		B string `plate:"b,path"`
		B string `plate:"b,form"`
	}
*/
func parse(v reflect.Value, jsonValue *fastjson.Value, jsonPath []string) error {
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.Type().NumField(); i++ {
			field := v.Type().Field(i)
			tag := field.Tag.Get(TagNameFetch)
			items := strings.Split(tag, TagOption)
			name := items[0]

			loc := "body"
			if len(items) > 1 {
				loc = items[1]
			}

			// json path
			newJsonPath := append(jsonPath, name)
			newJsonExist := jsonValue.Exists(newJsonPath...)
			if !newJsonExist {
				continue
			}

			switch loc {
			case "body":
				fmt.Println("json field:", field.Name)
				parseJson(v.Field(i), jsonValue, newJsonPath)

			case "header":
			case "path":
			case "form":
			default:
				return fmt.Errorf("")

			}
		}
	}

	return nil
}

func parseJson(v reflect.Value, jsonValue *fastjson.Value, jsonPath []string) error {
	vt := v.Type()
	switch vt.Kind() {
	case reflect.Struct:
		parse(v, jsonValue, jsonPath)
	case reflect.Slice, reflect.Array:

		arrayValue := jsonValue.GetArray(jsonPath...)
		for i := 0; i < len(arrayValue); i++ {
			arrayItem := reflect.New(vt.Elem())
			parse(arrayItem, jsonValue, nil)
			reflect.Append(v, arrayItem)
		}

	case reflect.Pointer:
		pointerType := vt.Elem()
		pointerValue := reflect.New(pointerType)
		parseJson(pointerValue.Elem(), jsonValue, jsonPath)
		v.Set(pointerValue)
	case reflect.Map:
		// get map key type and value type
		mapKeyType := v.Type().Key()
		mapValueType := v.Type().Elem()

		// check if map is nil then init it
		if v.IsNil() {
			v.Set(reflect.MakeMap(vt))
		}

		mapJsonValue := jsonValue.Get(jsonPath...)
		jsonOb, err := mapJsonValue.Object()
		if err != nil {
			return err
		}

		errStr := ""
		jsonOb.Visit(func(key []byte, childJsonValue *fastjson.Value) {
			mapKey, valueErr := getValueByType(key, mapKeyType)
			if valueErr != nil {
				errStr = errStr + fmt.Sprintf("get map key %s value type with error %v |", string(key), valueErr.Error())
				return
			}

			mapValueValue := reflect.New(mapValueType).Elem()
			parseJson(mapValueValue, childJsonValue, nil)
			v.SetMapIndex(mapKey, mapValueValue)
		})

		if len(errStr) > 0 {
			err = fmt.Errorf("%s", errStr)
		}

		return err

	case reflect.Bool:
		jv := jsonValue.GetBool(jsonPath...)
		v.SetBool(jv)
	case reflect.String:
		jv := jsonValue.GetStringBytes(jsonPath...)
		v.SetString(string(jv))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		jv := jsonValue.GetInt64(jsonPath...)
		fmt.Println("set int", jv, jsonPath)
		v.SetInt(jv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		jv := jsonValue.GetUint64(jsonPath...)
		v.SetUint(jv)
	case reflect.Float32, reflect.Float64:
		jv := jsonValue.GetFloat64(jsonPath...)
		v.SetFloat(jv)
	case reflect.Interface:
		// TODO

	}

	return nil
}

func getValueByType(key []byte, t reflect.Type) (reflect.Value, error) {
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
		// TODO
	}

	return res, nil
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

//func getMapKey(v reflect.Value) string {
//	map1 := make(map[reflect.Value]int)
//	switch v.Kind() {
//	case reflect.String:
//		return v.String()
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//		jv := jsonValue.GetInt64(jsonPath...)
//		v.SetInt(jv)
//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
//		jv := jsonValue.GetUint64(jsonPath...)
//		v.SetUint(jv)
//	case reflect.Float32, reflect.Float64:
//
//	}
//}
