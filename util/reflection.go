package util

import (
	"reflect"
	"strings"
)

func ObjectFields(t reflect.Type, ignores ...string) (fields []string) {
FIELDS:
	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		var tag = field.Tag.Get("db")
		if tag == "-" {
			continue
		}
		if tag == "" {
			fields = append(fields, strings.ToLower(field.Name))
		} else {
			for _, ignore := range ignores {
				if ignore == strings.Split(tag, ",")[0] {
					continue FIELDS
				}
			}
			fields = append(fields, strings.Split(tag, ",")[0])
		}
	}
	return fields
}

func ScanFields(val reflect.Value) []interface{} {
	v := val
	if val.Kind() == reflect.Ptr {
		v = val.Elem()
	}
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

func ValueFields(val reflect.Value, ignores ...string) []interface{} {
	v := val
	if val.Kind() == reflect.Ptr {
		v = val.Elem()
	}
	vs := []interface{}{}
FIELDS:
	for i := 0; i < v.NumField(); i++ {
		var tag = v.Type().Field(i).Tag.Get("db")
		if tag == "-" {
			continue
		}
		for _, ignore := range ignores {
			if ignore == strings.Split(tag, ",")[0] {
				continue FIELDS
			}
		}
		vs = append(vs, v.Field(i).Interface())
	}

	return vs
}

func IndexFieldTag(t reflect.Type, searchTag string) int {
	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		var tag = field.Tag.Get("db")
		if tag == searchTag {
			return i
		}
	}
	return -1
}
