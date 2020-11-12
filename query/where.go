package query

import "github.com/tessahoffmann/goqb/where"

func (self Query) Where(field string, operator string, value interface{}) Query {
	w := where.Where{
		Chain:    "AND",
		Field:    field,
		Operator: operator,
		Value:    value,
	}
	if len(self.WhereChain) == 0 {
		w.Chain = ""
	}

	self.WhereChain = append(self.WhereChain, w)

	return self
}

func (self Query) AndWhere(field string, operator string, value interface{}) Query {
	self.WhereChain = append(self.WhereChain, where.Where{
		Chain:    "AND",
		Field:    field,
		Operator: operator,
		Value:    value,
	})

	return self
}

func (self Query) AndWhereFunc(f func(Query) Query) Query {
	q := f(Query{})
	self.WhereChain = append(self.WhereChain, where.WhereGroup{
		Chain: "AND",
		Where: q.WhereChain,
	})
	return self
}

func (self Query) OrWhere(field string, operator string, value interface{}) Query {
	self.WhereChain = append(self.WhereChain, where.Where{
		Chain:    "OR",
		Field:    field,
		Operator: operator,
		Value:    value,
	})

	return self
}

func (self Query) OrWhereFunc(f func(Query) Query) Query {
	q := f(Query{})
	self.WhereChain = append(self.WhereChain, where.WhereGroup{
		Chain: "OR",
		Where: q.WhereChain,
	})
	return self
}
