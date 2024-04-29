package peak

type QueryBuilder interface {
	Build() BuildResult
}

type BuildResult struct {
	Sql    string
	Values map[string]any
}

type queryBuildersTree struct {
	filters       []*filterBuilder
	relationships []*relationshipBuilder
	selectors     []*selectorBuilder
	shapes        []*shapeBuilder
	temporaries   []*temporaryBuilder
	values        []*valuesBuilder
}

type queryPart struct {
	partType string
	name     string
	sql      string
	builder  QueryBuilder
	value    any
}

type Map = map[string]any
