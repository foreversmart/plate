package view

import (
	"reflect"
	"strings"
)

func fieldName(t reflect.StructField, tagOpt ...string) (res string, isValid bool) {
	res = t.Name

	if len(tagOpt) > 0 {
		tc := t.Tag.Get(tagOpt[0])
		items := strings.Split(tc, TagOption)
		res = items[0]
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
