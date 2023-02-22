# GoQB
GoQB is a simple Go Query Builder partly inspired by Laravel's Eloquent.

GoQB is instantiated like this:
```go
qb := goqb.NewGoQB(db) // db is of type *sql.DB and only supports MariaDB as of now
```

## Models
GoQB is based on models. Models are instantiated like this:
```go
type Book struct {
    ID        int    `db:"id"`
    Name      string `db:"name"`
    Author    string             // This will automatically be cast to be "author"
    Generated string `db:"-"`    // This will be ignored
}
var TableName = "books"

bookModel := qb.Model(TableName, Book{})
```
Tagged fields will automatically be used by GoQB to build the Query, untagged fields will be put into lower case (TODO: automatic camelcase to snake case conversion).

### Basic CRUD
Models can be directly queried.

**Get all models:**
```go
books := []Book{}
err := bookModel.All(&books)
```

**Get single model:**
```go
id := 1

book := Book{}
err := bookModel.Find(id, &book)
```

**Update model:**
```go
id := 1

type BookUpdate struct {
    Name string `db:"name"`
}
update := BookUpdate{Name: "New Name"}

err := qb.Model("books", update).Update(id, update)
```

**Create model:**
```go
insert := Book{
    Name:   "New Book",
    Author: "New Author",
}

err := qb.Model("books", insert).Create(&insert) // insert.ID will be filled with Auto-Increment data
```

### Where clauses
**Normal where clause:**
```go
books := []Book{}
bookModel.Where("author", "=", "Herman Melville").Get(&books)
```
Will result in:
```sql
SELECT id, name, author FROM books WHERE author = 'Herman Melville';
```

**Chaining where clauses:**
```go
books := []Book{}
bookModel.Where("author", "=", "Herman Melville").AndWhere("name", "LIKE", "%Moby%").Get(&books)
```
Will result in:
```sql
SELECT id, name, author FROM books WHERE author = 'Herman Melville' AND name LIKE '%Moby%';
```

**Nesting where clauses:**
```go
books := []Book{}
bookModel.Where("author", "=", "Herman Melville").AndWhereFunc(func(q Query) Query {
    return q.Where("name", "LIKE", "%Moby%").OrWhere("name", "=", "Israel Potter")
}).Get(&books)
```
Will result in:
```sql
SELECT id, name, author FROM books WHERE author = 'Herman Melville' AND (name LIKE '%Moby%' OR name = 'Israel Potter');
```

### Sorting
```go
books := []Book{}
bookModel.Where("author", "=", "Herman Melville").SortByDesc("name").Get(&books)
```
Will result in:
```sql
SELECT id, name, author FROM books WHERE author = 'Herman Melville' ORDER BY 'name' DESC;
```

## Relationships
Relationships are defined inside the main model's struct with the ```qb``` tag.
Only HasMany Relationships are implemented right now.

```go
type Library struct {
	ID   int    `db:"id"`
	Name string `db:"name"`

	Books []Book `qb:"hasMany"`
}
var TableName = "libraries"

libraryModel := qb.Model(TableName, Library{})
```

You can also specify the foreign table name and foreign key like this:
```go
Books []Book `qb:"hasMany,books,library_id"`
```
Otherwise this information will be automatically generated (only works with simple plurals and ```singular_id```-pattern foreign keys).

### Querying Relationships
```go
libraries := []Library{}
libraryModel.With("books").SortByDesc("name").Get(&libraries)
```
Will result in:
```sql
SELECT libraries.id, libraries.name, books.id, books.name, books.author FROM libraries JOIN books ON books.library_id = libraries.id;
```


## ToDo
* Automatically check fields on update and create, so you don't have to make a new model instance for those operations (currently it will throw an error because the insert/update fields will be generated from the struct you instantiated the model with)
* Add a ```Sort()``` function to Model