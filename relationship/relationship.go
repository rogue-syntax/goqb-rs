package relationship

import "fmt"

type Relationship interface {
	Table() string
	Fields() []string
	Type() string
	Index() int
	String(otherTable string, identifier string) string
}

type Relationships map[string]Relationship

type HasMany struct {
	TableName    string
	StructFields []string
	StructIndex  int
	ForeignID    string
}

func (self HasMany) Table() string {
	return self.TableName
}

func (self HasMany) Fields() []string {
	fields := []string{}
	for _, field := range self.StructFields {
		fields = append(fields, self.TableName+"."+field)
	}

	return fields
}

func (self HasMany) Type() string {
	return "hasMany"
}

func (self HasMany) Index() int {
	return self.StructIndex
}

func (self HasMany) String(otherTable string, identifier string) string {
	return fmt.Sprintf("JOIN %v ON %v = %v", self.TableName, self.TableName+"."+self.ForeignID, otherTable+"."+identifier)
}
