package fpgopg

import (
	"fmt"

	"github.com/cbrand/go-filterparams"
	"github.com/cbrand/go-filterparams/definition"
	"github.com/go-pg/pg/v10"
)

type Parser interface {
	AppendTo(*pg.Query, *filterparams.QueryData) *pg.Query
}

type defaultParser struct {
	naming Naming
}

func NewParser(naming Naming) Parser {
	return &defaultParser{
		naming: naming,
	}
}

func (pr defaultParser) AppendTo(q *pg.Query, queryData *filterparams.QueryData) *pg.Query {
	q = pr.appendFilter(q, queryData.GetFilter(), false)
	q = pr.appendOrder(q, queryData.GetOrders())
	return q
}

func (pr defaultParser) appendFilter(q *pg.Query, filter interface{}, or bool) *pg.Query {
	switch f := filter.(type) {
	case *definition.And:
		q = pr.appendAnd(q, f, or)
	case *definition.Or:
		q = pr.appendOr(q, f, or)
	case *definition.Negate:
		q = pr.appendNegate(q, f, or)
	case *definition.Parameter:
		q = pr.appendParameter(q, f, or)
	}

	return q
}

func (pr defaultParser) appendAnd(q *pg.Query, a *definition.And, or bool) *pg.Query {
	fn := func(q *pg.Query) (*pg.Query, error) {
		q = pr.appendFilter(q, a.Left, false)
		q = pr.appendFilter(q, a.Right, false)
		return q, nil
	}

	return pr.appendGroup(q, fn, or)
}

func (pr defaultParser) appendOr(q *pg.Query, o *definition.Or, or bool) *pg.Query {
	fn := func(q *pg.Query) (*pg.Query, error) {
		q = pr.appendFilter(q, o.Left, true)
		q = pr.appendFilter(q, o.Right, true)
		return q, nil
	}

	return pr.appendGroup(q, fn, or)
}

func (pr defaultParser) appendNegate(q *pg.Query, n *definition.Negate, or bool) *pg.Query {
	fn := func(q *pg.Query) (*pg.Query, error) {
		q = pr.appendFilter(q, n.Negated, false)
		return q, nil
	}

	return pr.appendGroup(q, fn, or)
}

func (pr defaultParser) appendGroup(q *pg.Query, fn func(q *pg.Query) (*pg.Query, error), or bool) *pg.Query {
	if or {
		q = q.WhereOrGroup(fn)
	} else {
		q = q.WhereGroup(fn)
	}

	return q
}

func (pr defaultParser) appendParameter(q *pg.Query, pm *definition.Parameter, or bool) *pg.Query {
	op := "="

	switch pm.Filter.Identification {
	// Default
	// case definition.FilterEq.Identification:
	// 	op = "="

	// IN seems a little strange for these filters. For now it is disabled.
	// case definition.FilterIn.Identification:
	// 	op = "IN"

	case definition.FilterLike.Identification:
		op = "LIKE"
	case definition.FilterILike.Identification:
		op = "ILIKE"
	case definition.FilterGt.Identification:
		op = ">"
	case definition.FilterGte.Identification:
		op = ">="
	case definition.FilterLt.Identification:
		op = "<"
	case definition.FilterLte.Identification:
		op = "<="
	}

	condition := fmt.Sprintf("? %s ?", op)
	ident := pg.Ident(pr.naming.Interpret(pm.Name))

	if or {
		q = q.WhereOr(condition, ident, pm.Value)
	} else {
		q = q.Where(condition, ident, pm.Value)
	}

	return q
}

func (pr defaultParser) appendOrder(q *pg.Query, orders []*definition.Order) *pg.Query {
	for _, order := range orders {
		ident := pg.Ident(pr.naming.Interpret(order.GetOrderBy()))

		if order.OrderDesc() {
			q = q.OrderExpr("? DESC", ident)
		} else {
			q = q.OrderExpr("? ASC", ident)
		}
	}

	return q
}
