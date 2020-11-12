package goqb

import (
	"database/sql"
	"reflect"
)

type GoQB struct {
	DB *sql.DB
}

func NewGoQB(db *sql.DB) *GoQB {
	return &GoQB{DB: db}
}

func (self GoQB) Close() {
	self.DB.Close()
}

func (self GoQB) Model(table string, obj interface{}) Model {
	t := reflect.TypeOf(obj)

	var fields = []string{"*"}

	if t != nil {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Struct {
			fields = ObjectFields(t)
		} else {
			if t.Kind() == reflect.Slice {
				f, ok := obj.([]string)
				if ok {
					fields = f
				}
			}
		}
	}

	return Model{
		Table:  table,
		Fields: fields,
		db:     self.DB,
	}
}
