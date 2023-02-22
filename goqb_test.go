package goqb

import (
	"log"
	"reflect"
	"testing"

	"github.com/rogue-syntax/goqb-rs/relationship"
	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	book := Book{}

	books := GoQB{nil}.Model("books", book)

	if !assert.Equal(t, books.String(), "SELECT id, name, author FROM books") {
		t.Error("SELECT query doesn't match")
	}
}

func TestModelWithHasMany(t *testing.T) {
	library := Library{}
	libraries := GoQB{nil}.Model("libraries", library)

	expect := relationship.Relationships{
		"books": relationship.HasMany{
			TableName:    "books",
			StructFields: []string{"id", "name", "author"},
			StructIndex:  2,
			ForeignID:    "library_id",
		},
	}

	log.Printf("libraries: %v", libraries.Relationships)
	log.Printf("expect: %v", expect)

	if !reflect.DeepEqual(libraries.Relationships, expect) {
		t.Error("Relationships don't match")
	}
}
