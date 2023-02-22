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
	return strings.Trim(fmt.Sprintf(`ORDER BY %s %s`, s.Offset, s.Limit), " ")
}
