package goqb

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/tessahoffmann/goqb/relationship"
	"github.com/tessahoffmann/goqb/util"
)

var (
	ErrNothingChanged = errors.New("nothing changed")
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
	var relations = relationship.Relationships{}

	if t != nil {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Struct {
			fields = util.ObjectFields(t)
			relations = util.ObjectRelations(table, t)
		} else {
			if t.Kind() == reflect.Slice {
				f, ok := obj.([]string)
				if ok {
					fields = f
				}
			}
		}
	}

	model := Model{
		Table:         table,
		Fields:        fields,
		Identifier:    "id",
		Relationships: relations,
		db:            self.DB,
	}

	return model
}
