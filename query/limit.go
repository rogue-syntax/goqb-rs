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
	return strings.Trim(fmt.Sprintf(`LIMIT %s OFFSET %s `, s.Limit, s.Offset), " ")
}

func (self Query) OffsetLimit(offset string, limit string) Query {
	self.Limit = Limit{Offset: offset, Limit: limit}
	return self
}
