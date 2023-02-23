package query

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rogue-syntax/goqb-rs/util"
	"github.com/rogue-syntax/goqb-rs/where"
)

type Fields []string

func (f Fields) String() string {
	return strings.Join(f, ", ")
}

type Query struct {
	Table      string
	Fields     Fields
	Identifier string
	Join       string
	JoinIndex  int
	WhereChain where.WhereChain
	Sort       Sort
	Limit      Limit
	DB         *sql.DB
}

func (self Query) Get(obj interface{}) error {
	if self.Join == "" {
		return self.getWithoutJoin(obj)
	} else {
		return self.getWithJoin(obj)
	}
}

func (self Query) getWithoutJoin(obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("obj has to be a pointer")
	}
	var s = v.Elem()
	if s.Kind() != reflect.Slice {
		return errors.New("obj has to be a pointer to a slice")
	}

	var query = self.String()
	rows, err := self.DB.Query(query, self.Args()...)
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

func (self Query) getWithJoin(obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("obj has to be a pointer")
	}
	var s = v.Elem()
	if s.Kind() != reflect.Slice {
		return errors.New("obj has to be a pointer to a slice")
	}
	var st = s.Type().Elem()

	var f = st.Field(self.JoinIndex).Type
	if f.Kind() != reflect.Slice {
		return errors.New("join field has to be a slice")
	}
	var ft = f.Elem()

	var query = self.String()
	rows, err := self.DB.Query(query, self.Args()...)
	if err != nil {
		return err
	}

	results := make(map[interface{}]reflect.Value)

	for rows.Next() {
		var inst = reflect.New(st)
		var foreign = reflect.New(ft)
		var scanFields = append(util.ScanFields(inst), util.ScanFields(foreign)...)
		err = rows.Scan(scanFields...)
		if err != nil {
			return err
		}

		i := inst.Elem().Field(util.IndexFieldTag(st, self.Identifier)).Interface()

		mapVal, ok := results[i]
		if ok {
			mapVal.Field(self.JoinIndex).Set(reflect.Append(mapVal.Field(self.JoinIndex), foreign.Elem()))
			results[i] = mapVal
		} else {
			elem := inst.Elem()
			elem.Field(self.JoinIndex).Set(reflect.MakeSlice(f, 0, 0))
			elem.Field(self.JoinIndex).Set(reflect.Append(elem.Field(self.JoinIndex), foreign.Elem()))
			results[i] = elem
		}
	}

	for _, val := range results {
		s.Set(reflect.Append(s, val))
	}

	return nil
}

func (self Query) String() string {
	operations := strings.Trim(strings.Join([]string{self.WhereChain.String(), self.Sort.String(), self.Limit.String()}, " "), " ")
	if self.Join != "" {
		operations = self.Join + " " + operations
	}
	return strings.Trim(fmt.Sprintf(`SELECT %s FROM %s %s`, self.Fields.String(), self.Table, operations), " ") + ";"
}

func (self Query) Args() []interface{} {
	return self.WhereChain.Args()
}
