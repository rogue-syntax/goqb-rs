package query

import (
	"github.com/tessahoffmann/goqb/where"
)

func (q Query) WhereJSONContains(field string, path string, value interface{}) Query {
	f := where.WhereFunction{
		Function: "JSON_CONTAINS",
		Field:    field,
		Flags:    []string{"'$'"},
		Value:    value,
	}
	if path != "" {
		f.Flags = []string{"'" + path + "'"}
	}

	q.WhereChain = q.WhereChain.Append("", f)
	return q
}
