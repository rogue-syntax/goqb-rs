package goqb

import (
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type Book struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Author    string
	Generated string `db:"-"`
}

func TestObjectFields(t *testing.T) {
	book := Book{}
	fields := ObjectFields(reflect.TypeOf(book))

	log.Printf("%v", fields)

	if !assert.ElementsMatch(t, fields, []string{"id", "name", "author"}) {
		t.Error("Elements don't match")
	}
}

func TestModel(t *testing.T) {
	book := Book{}

	books := GoQB{nil}.Model("books", book)

	if !assert.Equal(t, books.String(), "SELECT id, name, author FROM books;") {
		t.Error("SELECT query doesn't match")
	}
}

func TestModelAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, author FROM books;").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "author"}).AddRow(1, "Test Buch", "Max Mustermann").AddRow(2, "ABC Buch", "Maria Mustermann"))

	books := []Book{}

	err = GoQB{db}.Model("books", Book{}).All(&books)
	if err != nil {
		t.Error(err)
	}

	log.Printf("%v", books)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
