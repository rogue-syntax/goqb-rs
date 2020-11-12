package goqb

import (
	"reflect"
	"strings"
)

func ObjectFields(t reflect.Type) (fields []string) {
	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		var tag = field.Tag.Get("db")
		if tag == "-" {
			continue
		}
		if tag == "" {
			fields = append(fields, strings.ToLower(field.Name))
		} else {
			fields = append(fields, strings.Split(tag, ",")[0])
		}
	}
	return fields
}

func ScanFields(val reflect.Value) []interface{} {
	v := val.Elem()
	vs := []interface{}{}
	for i := 0; i < v.NumField(); i++ {
		var tag = v.Type().Field(i).Tag.Get("db")
		if tag == "-" {
			continue
		}
		vs = append(vs, v.Field(i).Addr().Interface())
	}

	return vs
}
