package view

import (
	"reflect"
	"strings"
)

// fieldName get filedName from reflect.StructField default name is field t's name
// full is ignore tagOpt select and filter, return default type's name
// full fetch from struct field tag split by : eg. "plate:"name,filter:full""
// tagOpt is where the field name fetch from the tag
// len(tagOpt) > 1 will filter the tag option
func fieldName(t reflect.StructField, full bool, tagOpt ...string) (res string, isFull bool, isValid bool) {
	res = t.Name
	isFull = full

	if len(tagOpt) > 0 {
		tc := t.Tag.Get(tagOpt[0])

		items := strings.Split(tc, TagOption)

		// when tc can't fetch default is t name
		if items[0] != "" {
			res = items[0]
		}

		// check is full
		fullItems := strings.Split(tc, TagConfig)
		if len(fullItems) > 1 && fullItems[1] == TagFull {
			isFull = true
		}

		if len(tagOpt) > 1 {

			if len(items) < 2 {
				isValid = false
			}

			// check tag option filter
			if len(items) > 1 && tagOpt[1] != items[1] {
				isValid = false
			}
		}

		// when field is fulled by parent all tag check is valid
		if full {
			isValid = true
		}

		isValid = isFieldNameValid(res)
	}

	return
}

func isFieldNameValid(n string) bool {
	if n == " " || n == "-" {
		return false
	}

	return true
}
