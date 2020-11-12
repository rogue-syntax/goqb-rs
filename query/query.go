package query

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tessahoffmann/goqb/util"
	"github.com/tessahoffmann/goqb/where"
)

type Fields []string

func (f Fields) String() string {
	return strings.Join(f, ", ")
}

type Query struct {
	Table      string
	Fields     Fields
	Identifier string
	WhereChain where.WhereChain
	Sort       Sort
	DB         *sql.DB
}

func (self Query) Get(obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("obj has to be a pointer")
	}
	var s = v.Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("obj has to be a pointer to a struct")
	}

	var query = self.String()
	rows, err := self.DB.Query(query, self.Args())
	if err != nil {
		return err
	}

	for rows.Next() {
		var inst = reflect.New(s.Type().Elem())
		err = rows.Scan(util.ScanFields(inst)...)
		if err != nil {
			return err
		}

		s.Set(reflect.Append(s, inst.Elem()))
	}

	return nil
}

func (self Query) String() string {
	whereSort := strings.Trim(strings.Join([]string{self.WhereChain.String(), self.Sort.String()}, " "), " ")
	return strings.Trim(fmt.Sprintf(`SELECT %s FROM %s %s`, self.Fields.String(), self.Table, whereSort), " ") + ";"
}

func (self Query) Args() interface{} {
	return self.WhereChain.Args()
}
