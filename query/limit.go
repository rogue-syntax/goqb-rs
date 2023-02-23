package query

import (
	"fmt"
	"strings"
)

type Limit struct {
	Offset string
	Limit  string
}

func (s Limit) String() string {
	if s.Offset == "" && s.Limit == "" {
		return ""
	}
	return strings.Trim(fmt.Sprintf(`OFFSET %s LIMIT %s`, s.Offset, s.Limit), " ")
}

func (self Query) OffestLimit(offset string, limit string) Query {
	self.Limit = Limit{Offset: offset, Limit: limit}
	return self
}
