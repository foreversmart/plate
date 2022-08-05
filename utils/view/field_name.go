package view

import (
	"reflect"
	"strings"
)

// fieldName get filedName from reflect.StructField default name is field t's name
// tagOpt is where the field name fetch from the tag
// len(tagOpt) > 1 will filter the tag option
func fieldName(t reflect.StructField, tagOpt ...string) (res string, isValid bool) {
	res = t.Name

	if len(tagOpt) > 0 {
		tc := t.Tag.Get(tagOpt[0])
		items := strings.Split(tc, TagOption)
		res = items[0]

		if len(tagOpt) > 1 {

			if len(items) < 2 {
				return "", false
			}

			// check tag option filter
			if tagOpt[1] != items[1] {
				return "", false
			}
		}
	}

	isValid = isFieldNameValid(res)

	return
}

func isFieldNameValid(n string) bool {
	if n == " " || n == "-" {
		return false
	}

	return true
}
