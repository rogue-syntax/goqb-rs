package goqb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Model struct {
	Table  string
	Fields []string
	db     *sql.DB
}

func (self Model) All(obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("obj has to be a pointer")
	}
	var s = v.Elem()
	if s.Kind() != reflect.Slice {
		return errors.New("obj has to be a pointer to a slice")
	}

	var query = self.String()
	rows, err := self.db.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var inst = reflect.New(s.Type().Elem())
		err = rows.Scan(ScanFields(inst)...)
		if err != nil {
			return err
		}

		s.Set(reflect.Append(s, inst.Elem()))
	}

	return nil
}

func (self Model) String() string {
	return fmt.Sprintf("SELECT %s FROM %s;", strings.Join(self.Fields, ", "), self.Table)
}
