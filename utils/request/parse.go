package request

import (
	"encoding/json"
	"fmt"
	"github.com/foreversmart/plate/utils/view"
	"reflect"
	"strings"

	"github.com/foreversmart/plate/utils/val"
	"github.com/valyala/fastjson"
)

// ParseRequest parse req Requester into a value v
// v must be struct value
func ParseRequest(req Requester, v reflect.Value) error {
	jsonValue, meta, err := fetchJsonAndMeta(req)
	if err != nil {
		return err
	}

	return parse(v, jsonValue, meta, nil)
}

const (
	TagNameFetch   = "plate"
	TagNameCheck   = "check"
	TagOption      = ","
	TagOptionSplit = ":"

	TagOptionInline = "inline"

	// loc in request
	LocHeader = "header"
	LocBody   = "body"
	LocPath   = "path"
	LocForm   = "form"
	LocQuery  = "query"
	LocMid    = "mid"

	// loc in resp
	LocResp = "resp"
)

/*
	 Parse: parse a json value and meta to struct reflect value

		type A struct {
			B string `plate:"b,body" check:"int>10"`
			B string `plate:"b,header"`
			B string `plate:"b,path"`
			B string `plate:"b,form"`
			B string `plate:"b,mid:inline"`
		}

		type Mid struct {
			B string `mid:"b"`
		}

//
// meta must be golang basic type such as int, string, bool , float
*/
func Parse(v reflect.Value, jsonValue *fastjson.Value, meta map[string]map[string][]string, mid *view.View, jsonPath []string) (err error) {
	err = parse(v, jsonValue, meta, jsonPath)
	if err != nil {
		return err
	}

	// set middleware value
	if mid != nil {
		err = mid.SetStructValue(v, true, false, "plate", LocMid)
	}

	return
}

/*
parse: parse a json value and meta to struct reflect value

	type A struct {
		B string `plate:"b,body" check:"int>10"`
		B string `plate:"b,header"`
		B string `plate:"b,path"`
		B string `plate:"b,form"`
		B string `plate:"b,mid"`
		B string `plate:"b,body:inline"`
	}

	type Mid struct {
		B string `mid:"b"`
	}

//
// meta must be golang basic type such as int, string, bool , float
// v must struct or the parse is invalid
*/
func parse(v reflect.Value, jsonValue *fastjson.Value, meta map[string]map[string][]string, jsonPath []string) error {
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.Type().NumField(); i++ {
			field := v.Type().Field(i)
			tag := field.Tag.Get(TagNameFetch)
			items := strings.Split(tag, TagOptionSplit)
			option := ""
			if len(items) > 1 {
				option = items[1]
			}

			items = strings.Split(items[0], TagOption)
			name := items[0]

			loc := LocBody
			if len(items) > 1 {
				loc = items[1]
			}

			switch loc {
			case LocBody:
				if option == TagOptionInline {
					err := parseJson(v.Field(i), jsonValue, meta, jsonPath)
					if err != nil {
						fmt.Printf("parseJson %s with error %v \n", jsonPath, err)
						continue
					}
					continue
				}

				// json path
				newJsonPath := append(jsonPath, name)
				newJsonExist := jsonValue.Exists(newJsonPath...)
				if !newJsonExist {
					continue
				}

				err := parseJson(v.Field(i), jsonValue, meta, newJsonPath)
				if err != nil {
					fmt.Printf("parseJson %s with error %v \n", newJsonPath, err)
					continue
				}
			default:
				if option == TagOptionInline {
					vf := v.Field(i)
					for vf.Kind() == reflect.Ptr {
						if vf.IsNil() {
							//return nil, fmt.Errorf("cant fetch view from a nil object %s", v.Type().String())
							ot := vf.Type().Elem()
							vf.Set(reflect.New(ot))
						}

						child := vf.Elem()
						// break loop pointer
						if vf == child {
							break
						}

						vf = child
					}

					err := parse(vf, jsonValue, meta, jsonPath)
					if err != nil {
						return err
					}
					continue
				}

				metaMap, ok := meta[loc]
				if !ok {
					continue
				}

				metaValue, ok := metaMap[name]
				if !ok {
					continue
				}

				// default parse the first value
				err := val.SetValueByString(val.SettableValue(v.Field(i)), metaValue[0])
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

var (
	marshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
)

func parseJson(v reflect.Value, jsonValue *fastjson.Value, meta map[string]map[string][]string, jsonPath []string) error {
	// handle Unmarshaler case
	if v.Kind() != reflect.Pointer && v.Type().Name() != "" && v.CanAddr() {
		vp := v.Addr()
		vi, ok := vp.Interface().(json.Unmarshaler)
		if ok {
			return vi.UnmarshalJSON(jsonValue.Get(jsonPath...).MarshalTo(nil))
		}
	}

	vi, ok := v.Interface().(json.Unmarshaler)
	if ok {
		return vi.UnmarshalJSON(jsonValue.Get(jsonPath...).GetStringBytes())
	}

	switch v.Kind() {
	case reflect.Struct:
		parse(v, jsonValue, meta, jsonPath)
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
			mapKey, valueErr := val.NewValueByBytes(key, mapKeyType)
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
