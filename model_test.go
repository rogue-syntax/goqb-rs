package goqb

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type Library struct {
	ID   int    `db:"id"`
	Name string `db:"name"`

	Books []Book `qb:"hasMany"`
}

type Book struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Author    string
	Generated string `db:"-"`

	LibraryID int `db:"-"`
}

func TestModelFind(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, author FROM books WHERE id = \\?;").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "author"}).AddRow(1, "Test Buch", "Max Mustermann").AddRow(2, "ABC Buch", "Maria Mustermann"))

	book := Book{}

	err = GoQB{db}.Model("books", Book{}).Find(1, &book)
	if err != nil {
		t.Error(err)
	}

	log.Printf("%v", book)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModelFindErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, author FROM books WHERE id = \\?;").WillReturnError(sql.ErrNoRows)

	book := Book{}

	err = GoQB{db}.Model("books", Book{}).Find(3, &book)
	if err != sql.ErrNoRows {
		t.Error("Error was not ErrNoRows")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
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

func TestModelUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE books SET name = \\?, author = \\? WHERE id = \\?;").WithArgs("newName", "newAuthor", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	type Update struct {
		Name   string
		Author string
	}

	update := Update{"newName", "newAuthor"}

	err = GoQB{db}.Model("books", update).Update(1, update)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModelCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO books").WithArgs("Testbuch", "Testautor").WillReturnResult(sqlmock.NewResult(15, 1))

	insert := Book{
		Name:   "Testbuch",
		Author: "Testautor",
	}

	err = GoQB{db}.Model("books", insert).Create(&insert)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if insert.ID != 15 {
		t.Error("Expected ID to be 15, was not")
	}
}

func TestModelDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM books").WithArgs(15).WillReturnResult(sqlmock.NewResult(0, 1))

	err = GoQB{db}.Model("books", Book{}).Delete(15)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModelDeleteNoResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM books").WithArgs(15).WillReturnResult(sqlmock.NewResult(0, 0))

	err = GoQB{db}.Model("books", Book{}).Delete(15)
	if err != ErrNothingChanged {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModelWith(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var expect = []Library{
		{
			ID:   1,
			Name: "Barack Obama Presidential Library",
			Books: []Book{
				{
					ID:     1,
					Name:   "Test Buch",
					Author: "Max Mustermann",
				},
				{
					ID:     2,
					Name:   "ABC Buch",
					Author: "Maria Mustermann",
				},
			},
		},
	}

	mock.ExpectQuery("SELECT libraries.id, libraries.name, books.id, books.name, books.author FROM libraries JOIN books ON books.library_id = libraries.id").WillReturnRows(
		sqlmock.NewRows(
			[]string{"libraries.id", "libraries.name", "books.id", "books.name", "books.author"}).AddRow(
			1, "Barack Obama Presidential Library", 1, "Test Buch", "Max Mustermann").AddRow(
			1, "Barack Obama Presidential Library", 2, "ABC Buch", "Maria Mustermann"))

	var libraries = []Library{}
	err = GoQB{db}.Model("libraries", Library{}).With("books").Get(&libraries)
	if err != nil {
		t.Error(err)
	}

	log.Printf("%v", libraries)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(libraries, expect) {
		t.Error("results don't match")
	}
}
