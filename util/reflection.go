package util

import (
	"reflect"
	"strings"

	"github.com/tessahoffmann/goqb/relationship"
)

func ObjectFields(t reflect.Type, ignores ...string) (fields []string) {
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
FIELDS:
	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		if field.Tag.Get("qb") != "" {
			continue
		}

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

func ObjectRelations(table string, t reflect.Type) relationship.Relationships {
	var relations = make(map[string]relationship.Relationship)

	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		var tag = field.Tag.Get("qb")
		if tag == "-" || tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		if len(parts) == 0 {
			continue
		}
		if len(parts) > 0 {
			if parts[0] == "hasMany" {
				var rel = relationship.HasMany{
					TableName:    strings.ToLower(field.Name),
					ForeignID:    singularize(table) + "_id",
					StructFields: ObjectFields(field.Type),
					StructIndex:  i,
				}
				if len(parts) > 1 {
					rel.TableName = parts[1]
				}
				if len(parts) > 2 {
					rel.ForeignID = parts[2]
				}
				relations[strings.ToLower(field.Name)] = rel
			}
		}
	}
	return relations
}

func ScanFields(val reflect.Value) []interface{} {
	v := val
	if val.Kind() == reflect.Ptr {
		v = val.Elem()
	}
	vs := []interface{}{}
	for i := 0; i < v.NumField(); i++ {
		var field = v.Type().Field(i)
		if field.Tag.Get("qb") != "" {
			continue
		}

		var tag = field.Tag.Get("db")
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

func singularize(s string) string {
	if strings.HasSuffix(s, "ies") {
		return strings.TrimRight(s, "ies") + "y"
	}
	if strings.HasSuffix(s, "es") {
		return strings.TrimRight(s, "es")
	}
	if strings.HasSuffix(s, "i") {
		return strings.TrimRight(s, "i") + "us"
	}
	if strings.HasSuffix(s, "a") {
		return strings.TrimRight(s, "a") + "on"
	}
	return s
}
