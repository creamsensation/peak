package peak

type TemporaryBuilder interface {
	QueryBuilder
}

type temporaryBuilder struct {
	builder QueryBuilder
	name    string
}

func Temporary(name string, qb QueryBuilder) TemporaryBuilder {
	return &temporaryBuilder{
		builder: qb,
		name:    name,
	}
}

func (b *temporaryBuilder) Build() BuildResult {
	build := b.builder.Build()
	return BuildResult{
		Sql:    "(" + build.Sql + ") AS " + b.name,
		Values: build.Values,
	}
}
