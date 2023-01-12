package common

import (
	"crud_service/app/repos"
	"reflect"
	"strings"
)

func CheckRequiredFields(d repos.Data) bool {
	fields := reflect.ValueOf(d).Elem()
	for i := 0; i < fields.NumField(); i++ {
		tag := fields.Type().Field(i).Tag.Get("validate")
		if strings.Contains(tag, "required") && fields.Field(i).IsZero() {
			return false
		}
	}
	return true
}
