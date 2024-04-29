package peak

type RemoveRepository[R result] interface {
	QueryBuilder
	Run(runner ...Runner) (R, error)
	MustRun(runner ...Runner) R
}

type removeRepository[E entity, R result] struct {
	*repository[E, R]
	filters   []*filterBuilder
	selectors []*selectorBuilder
}

func (r *removeRepository[E, R]) Build() BuildResult {
	values := make(map[string]any)
	selectorsExist := len(r.selectors) > 0
	e := any(r.entity).(entity)
	q := createSqlBuilder().
		Q("DELETE").
		Q("FROM " + e.Table()).
		Q("AS " + e.Alias())
	
	// Where
	buildBeforeAggregationFilters(q, r.filters, &values)
	
	q.If(!selectorsExist, "RETURNING *").
		If(selectorsExist, "RETURNING "+buildFieldsSql(r.selectors...))
	
	return BuildResult{q.Build(), values}
}

func (r *removeRepository[E, R]) Run(runner ...Runner) (R, error) {
	res := new(R)
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

func (r *removeRepository[E, R]) MustRun(runner ...Runner) R {
	res, err := r.Run(runner...)
	if err != nil {
		panic(err)
	}
	return res
}
