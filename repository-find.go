package peak

import (
	"reflect"
)

type FindRepository[R result] interface {
	QueryBuilder
	Run(runner ...Runner) (R, error)
	MustRun(runner ...Runner) R
}

type findRepository[E entity, R result] struct {
	*repository[E, R]
	filters       []*filterBuilder
	relationships []*relationshipBuilder
	selectors     []*selectorBuilder
	shapes        []*shapeBuilder
}

func (r *findRepository[E, R]) Build() BuildResult {
	var fieldsSql string
	values := make(map[string]any)
	e := any(r.entity).(entity)
	fields := e.Fields()
	selectorsExist := len(r.selectors) > 0
	if !selectorsExist {
		fieldsSql = buildFieldsSql(fields...)
	}
	if selectorsExist {
		fieldsSql = buildFieldsSql(r.selectors...)
	}
	q := createSqlBuilder().
		Q("SELECT").
		If(doesExistDistinct(r.shapes), "DISTINCT").
		Q(fieldsSql).
		Q("FROM " + e.Table()).
		Q("AS " + e.Alias())
	
	// Joins
	buildJoins(q, r.relationships, fields)
	
	// Where
	buildBeforeAggregationFilters(q, r.filters, &values)
	
	// Group shapes
	groupShapes := buildGroupShapes(r.shapes)
	q = q.If(len(groupShapes) > 0, groupShapes)
	
	// Having
	buildAfterAggregationFilters(q, r.filters, &values)
	
	// Order, Limit, Offset
	nonGroupShapes := buildNonGroupShapes(r.shapes)
	q = q.If(len(nonGroupShapes) > 0, nonGroupShapes)
	
	return BuildResult{q.Build(), values}
}

func (r *findRepository[E, R]) Run(runner ...Runner) (R, error) {
	res := new(R)
	t := reflect.TypeOf(res)
	if t.Elem().Kind() == reflect.Slice {
		slicePtr := reflect.New(t.Elem())
		slicePtr.Elem().Set(reflect.MakeSlice(t.Elem(), 0, 0))
		*res = slicePtr.Elem().Interface().(R)
	}
	if r.db == nil {
		return *res, ErrorMissingDatabase
	}
	if len(runner) > 0 {
		return *res, nil
	}
	b := r.Build()
	if err := r.db.Q(b.Sql, b.Values).Exec(res); err != nil {
		return *res, err
	}
	return *res, nil
}

func (r *findRepository[E, R]) MustRun(runner ...Runner) R {
	res, err := r.Run(runner...)
	if err != nil {
		panic(err)
	}
	return res
}
