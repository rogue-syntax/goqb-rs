package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	q := Query{
		Fields: []string{"*"},
		Table:  "testing",
	}
	q = q.SortByDesc("sort_col")

	if !assert.Equal(t, q.String(), "SELECT * FROM testing ORDER BY sort_col DESC;") {
		t.Error("Query string was incorrect")
	}
}

func TestWhere(t *testing.T) {
	q := Query{
		Fields: []string{"*"},
		Table:  "testing",
	}
	q = q.Where("column1", "=", 1).OrWhereFunc(func(q Query) Query {
		return q.Where("column2", ">", "2").AndWhere("column3", "<", 3)
	})

	if !assert.Equal(t, q.String(), "SELECT * FROM testing WHERE column1 = ? OR (column2 > ? AND column3 < ?);") {
		t.Error("Query string was incorrect")
	}

	if !assert.ElementsMatch(t, q.Args(), []interface{}{1, "2", 3}) {
		t.Error("Args were incorrect")
	}
}

func TestWhereJSONContains(t *testing.T) {
	q := Query{
		Fields: []string{"*"},
		Table:  "testing",
	}
	q = q.WhereJSONContains("json_field", "$", []int{1})

	if !assert.Equal(t, q.String(), "SELECT * FROM testing WHERE JSON_CONTAINS(json_field, ?, '$');") {
		t.Error("Query string was incorrect")
	}

	if !assert.ElementsMatch(t, q.Args(), []interface{}{[]int{1}}) {
		t.Error("Args were incorrect")
	}
}

func TestWhereHasMany(t *testing.T) {
	q := Query{
		Fields: []string{"*"},
		Table:  "books",
	}
	q = q.WhereHasMany("id", "IN", 1, "library_books", "book_id", "library_id")

	if !assert.Equal(t, q.String(), "SELECT * FROM books WHERE id IN (SELECT library_books.book_id FROM library_books WHERE library_books.library_id = ?);") {
		t.Error("Query string was incorrect")
	}

	if !assert.ElementsMatch(t, q.Args(), []interface{}{1}) {
		t.Error("Args were incorrect")
	}
}
