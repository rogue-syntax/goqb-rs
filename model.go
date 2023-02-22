package goqb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tessahoffmann/goqb/query"
	"github.com/tessahoffmann/goqb/relationship"
	"github.com/tessahoffmann/goqb/util"
	"github.com/tessahoffmann/goqb/where"
)

type Model struct {
	Table         string
	Fields        []string
	Identifier    string
	Relationships relationship.Relationships
	db            *sql.DB
}

func (self Model) With(relation string) query.Query {
	rel, ok := self.Relationships[relation]
	if !ok {
		return query.Query{
			Table:      self.Table,
			Fields:     self.Fields,
			Identifier: self.Identifier,
			DB:         self.db,
		}
	}

	fields := []string{}

	for _, field := range self.Fields {
		fields = append(fields, self.Table+"."+field)
	}

	return query.Query{
		Table:      self.Table,
		Fields:     append(fields, rel.Fields()...),
		Identifier: self.Identifier,
		Join:       rel.String(self.Table, self.Identifier),
		JoinIndex:  rel.Index(),
		DB:         self.db,
	}
}

func (self Model) Query() query.Query {
	return query.Query{
		Table:      self.Table,
		Fields:     self.Fields,
		Identifier: self.Identifier,
		DB:         self.db,
	}
}

func (self Model) Where(field string, operator string, value interface{}) query.Query {
	return query.Query{
		Table:      self.Table,
		Fields:     self.Fields,
		Identifier: self.Identifier,
		WhereChain: []where.IWhere{
			where.Where{
				Field:    field,
				Operator: operator,
				Value:    value,
			},
		},
		DB: self.db,
	}
}

func (self Model) Find(id interface{}, obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("obj has to be a pointer")
	}
	var s = v.Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("obj has to be a pointer to a struct")
	}

	var query = fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?;`, strings.Join(self.Fields, ", "), self.Table, self.Identifier)
	err := self.db.QueryRow(query, id).Scan(util.ScanFields(s)...)
	if err != nil {
		return err
	}

	return nil
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
		err = rows.Scan(util.ScanFields(inst)...)
		if err != nil {
			return err
		}

		s.Set(reflect.Append(s, inst.Elem()))
	}

	return nil
}

func (self Model) Update(id interface{}, obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		return errors.New("obj has to be a struct")
	}

	vals := append(util.ValueFields(v), id)

	var query = fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?;`, self.Table, strings.Join(self.Fields, " = ?, "), self.Identifier)
	_, err := self.db.Exec(query, vals...)
	if err != nil {
		return err
	}

	return nil
}

func (self Model) Create(obj interface{}) error {
	var v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("obj has to be a pointer")
	}
	var s = v.Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("obj has to be a pointer to a struct")
	}

	vals := util.ValueFields(v, self.Identifier)

	placeholder := strings.TrimRight(strings.Repeat("?, ", len(vals)), ", ")

	var query = fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s);`, self.Table, strings.Join(util.ObjectFields(s.Type(), self.Identifier), ", "), placeholder)
	result, err := self.db.Exec(query, vals...)
	if err != nil {
		return err
	}

	liid, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.Field(util.IndexFieldTag(s.Type(), self.Identifier)).SetInt(liid)

	return nil
}

func (self Model) Delete(id interface{}) error {
	var query = fmt.Sprintf(`DELETE FROM %s WHERE %s = ?;`, self.Table, self.Identifier)
	result, err := self.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNothingChanged
	}
	return nil
}

func (self Model) String() string {
	return fmt.Sprintf("SELECT %s FROM %s;", strings.Join(self.Fields, ", "), self.Table)
}
