package where

import (
	"fmt"
	"strings"
)

type IWhere interface {
	SetChain(chain string)
	WhereString() string
	Args() []interface{}
}

type Where struct {
	Chain    string // either empty, AND or OR
	Field    string
	Operator string
	Value    interface{}
}

func (w Where) SetChain(chain string) {
	w.Chain = chain
}

func (w Where) WhereString() string {
	return strings.Trim(fmt.Sprintf("%s %s %s ?", w.Chain, w.Field, w.Operator), " ")
}

func (w Where) Args() []interface{} {
	return []interface{}{w.Value}
}

type WhereFunction struct {
	Chain    string
	Field    string
	Flags    []string
	Function string
	Value    interface{}
}

func (wf WhereFunction) SetChain(chain string) {
	wf.Chain = chain
}

func (wf WhereFunction) WhereString() string {
	return strings.Trim(fmt.Sprintf("%s %s(%s, %s)", wf.Chain, wf.Function, wf.Field, strings.Join(append([]string{"?"}, wf.Flags...), ", ")), " ")
}

func (wf WhereFunction) Args() []interface{} {
	return []interface{}{wf.Value}
}

type WhereSubquery struct {
	Chain          string
	Field          string
	Operator       string
	PivotTable     string
	IDField        string
	ForeignIDField string
	Value          interface{}
}

func (ws WhereSubquery) SetChain(chain string) {
	ws.Chain = chain
}

func (ws WhereSubquery) WhereString() string {
	return strings.Trim(fmt.Sprintf("%s %s %s (SELECT %s.%s FROM %s WHERE %s.%s = ?)", ws.Chain, ws.Field, ws.Operator, ws.PivotTable, ws.IDField, ws.PivotTable, ws.PivotTable, ws.ForeignIDField), " ")
}

func (ws WhereSubquery) Args() []interface{} {
	return []interface{}{ws.Value}
}

type WhereGroup struct {
	Chain string
	Where []IWhere
}

func (wg WhereGroup) SetChain(chain string) {
	wg.Chain = chain
}

func (wg WhereGroup) WhereString() string {
	var chain = []string{}

	for _, where := range wg.Where {
		chain = append(chain, where.WhereString())
	}

	return fmt.Sprintf("%s (%s)", wg.Chain, strings.Join(chain, " "))
}

func (wg WhereGroup) Args() []interface{} {
	var chain = []interface{}{}

	for _, where := range wg.Where {
		chain = append(chain, where.Args()...)
	}

	return chain
}

type WhereChain []IWhere

func (wc WhereChain) String() string {
	if len(wc) == 0 {
		return ""
	}

	var chain = []string{}
	for _, where := range wc {
		chain = append(chain, where.WhereString())
	}

	return "WHERE " + strings.Join(chain, " ")
}

func (wc WhereChain) Args() []interface{} {
	var chain = []interface{}{}

	for _, where := range wc {
		chain = append(chain, where.Args()...)
	}

	return chain
}

func (wc WhereChain) Append(chain string, w IWhere) WhereChain {
	if len(wc) == 0 {
		w.SetChain("")
	} else {
		if chain == "OR" {
			w.SetChain("OR")
		} else {
			w.SetChain("AND")
		}
	}

	wc = append(wc, w)
	return wc
}
