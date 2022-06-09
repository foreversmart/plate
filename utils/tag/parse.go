package tagger

import (
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Requester interface {
	Request() *http.Request
	Params() map[string][]string // path param value
}

type ParseRequester interface {
}

func ParseRequest(req Requester, v reflect.Value) error {
	body, err := ioutil.ReadAll(req.Request().Body)
	defer req.Request().Body.Close()
	if err != nil {
		return err
	}

	jsonValue, err := fastjson.Parse(string(body))
	if err != nil {
		return err
	}

	meta := make(map[string]map[string][]string)
	meta[LocHeader] = req.Request().Header

	formErr := req.Request().ParseForm()
	// indicate is form request
	if formErr == nil {
		meta[LocForm] = req.Request().PostForm
	}
	meta[LocQuery] = req.Request().URL.Query()
	meta[LocPath] = req.Params()

	return Parse(v, jsonValue, nil, nil)
}

const (
	TagNameFetch = "plate"
	TagNameCheck = "check"
	TagOption    = ","

	LocHeader = "header"
	LocBody   = "body"
	LocPath   = "path"
	LocForm   = "form"
	LocQuery  = "query"
)

/*
	Parse: parse a json value and meta to struct reflect value
	type A struct {
		B string `plate:"b,body" check:"int>10"`
		B string `plate:"b,header"`
		B string `plate:"b,path"`
		B string `plate:"b,form"`
	}
	//
	// meta must be golang basic type such as int, string, bool , float
*/
func Parse(v reflect.Value, jsonValue *fastjson.Value, meta map[string]map[string]string, jsonPath []string) error {
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
				parseJson(v.Field(i), jsonValue, meta, newJsonPath)
			default:
				metaMap, ok := meta[loc]
				if !ok {
					continue
				}

				metaValue, ok := metaMap[name]
				if !ok {
					continue
				}

				parseValueByText(v.Field(i), metaValue)

			}
		}
	}

	return nil
}

func parseJson(v reflect.Value, jsonValue *fastjson.Value, meta map[string]map[string]string, jsonPath []string) error {
	switch v.Kind() {
	case reflect.Struct:
		Parse(v, jsonValue, meta, jsonPath)
	case reflect.Slice, reflect.Array:
		vt := v.Type()
		// fastjson array value
		arrayValue := jsonValue.GetArray(jsonPath...)

		// if array value is nil or length = 0 return
		if len(arrayValue) == 0 {
			return nil
		}

		// slice is nil or length is nil init with len
		if v.IsNil() || v.Len() == 0 {
			newArray := reflect.MakeSlice(vt, len(arrayValue), len(arrayValue))
			v.Set(newArray)
		}

		// parse array item
		for i := 0; i < len(arrayValue); i++ {
			parseJson(v.Index(i), arrayValue[i], meta, nil)
		}

	case reflect.Pointer:
		vt := v.Type()
		pointerType := vt.Elem()
		pointerValue := reflect.New(pointerType)
		parseJson(pointerValue.Elem(), jsonValue, meta, jsonPath)
		v.Set(pointerValue)
	case reflect.Map:
		vt := v.Type()

		// get map key type and value type
		mapKeyType := v.Type().Key()
		mapValueType := v.Type().Elem()

		// check if map is nil then init it
		if v.IsNil() {
			v.Set(reflect.MakeMap(vt))
		}

		// get fastjson map object for visit
		mapJsonValue := jsonValue.Get(jsonPath...)
		jsonOb, err := mapJsonValue.Object()
		if err != nil {
			return err
		}

		errStr := ""
		// visit map each key and value and parse
		jsonOb.Visit(func(key []byte, childJsonValue *fastjson.Value) {
			mapKey, valueErr := getValueByType(key, mapKeyType)
			if valueErr != nil {
				errStr = errStr + fmt.Sprintf("get map key %s value type with error %v |", string(key), valueErr.Error())
				return
			}

			mapValueValue := reflect.New(mapValueType).Elem()
			err = parseJson(mapValueValue, childJsonValue, meta, nil)
			if err != nil {
				errStr = errStr + " | " + err.Error() + " | "
			}
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
		v.SetInt(jv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		jv := jsonValue.GetUint64(jsonPath...)
		v.SetUint(jv)
	case reflect.Float32, reflect.Float64:
		jv := jsonValue.GetFloat64(jsonPath...)
		v.SetFloat(jv)
	case reflect.Interface:
		jv := jsonValue.Get(jsonPath...)
		switch jv.Type() {
		case fastjson.TypeNull:
			// no need to do
		case fastjson.TypeObject:
			// Golang json unmarshal to map[string]interface{} so do the same
			a := make(map[string]interface{})
			parseJson(reflect.ValueOf(&a).Elem(), jsonValue, meta, jsonPath)
			v.Set(reflect.ValueOf(a))
		case fastjson.TypeArray:
			jvArray := jv.GetArray()
			a := make([]interface{}, len(jv.GetArray()))
			for i := 0; i < len(jvArray); i++ {
				parseJson(reflect.ValueOf(&a[i]).Elem(), jvArray[i], meta, nil)
			}
			v.Set(reflect.ValueOf(a))
		case fastjson.TypeString:
			b := jv.GetStringBytes()
			v.Set(reflect.ValueOf(string(b)))
		case fastjson.TypeNumber:
			b := jv.GetFloat64()
			v.Set(reflect.ValueOf(b))
		case fastjson.TypeTrue, fastjson.TypeFalse:
			b := jv.GetBool()
			v.Set(reflect.ValueOf(b))
		}
	default:
		return fmt.Errorf("parse json type is invalid %s", v.String())

	}

	return nil
}

func parseValueByText(value reflect.Value, s string) error {
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
		// TODO need dynamic check the type but json never happen
		// default support as string
		res.Set(reflect.ValueOf(string(key)))
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
