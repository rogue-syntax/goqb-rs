package query

import (
	"fmt"
	"strings"
)

type Sort struct {
	Field     string
	Direction string
}

func (s Sort) String() string {
	if s.Field == "" && s.Direction == "" {
		return ""
	}
	return strings.Trim(fmt.Sprintf(`SORT BY %s %s`, s.Field, s.Direction), " ")
}

func (self Query) SortBy(field string) Query {
	self.Sort = Sort{Field: field, Direction: "ASC"}
	return self
}

func (self Query) SortByDesc(field string) Query {
	self.Sort = Sort{Field: field, Direction: "DESC"}
	return self
}
