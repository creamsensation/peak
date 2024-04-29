package peak

import (
	"fmt"
	"strings"
)

type SelectorBuilder interface {
	QueryBuilder
	Group(qb QueryBuilder) SelectorBuilder
	Field(field QueryBuilder) SelectorBuilder
	As(value string) SelectorBuilder
	Count(f QueryBuilder) SelectorBuilder
	Sum(f QueryBuilder) SelectorBuilder
	Min(qb QueryBuilder) SelectorBuilder
	Max(qb QueryBuilder) SelectorBuilder
	Avg(qb QueryBuilder) SelectorBuilder
}

type selectorBuilder struct {
	parts      []queryPart
	fieldsOnly bool
}

const (
	selectorGroupPart     = "selector-group"
	selectorFieldPart     = "selector-field"
	selectorAliasPart     = "selector-alias"
	selectorAggregatePart = "selector-aggregate"
)

func Selector(fields ...QueryBuilder) SelectorBuilder {
	b := &selectorBuilder{
		parts: make([]queryPart, 0),
	}
	if len(fields) > 0 {
		for _, f := range fields {
			b.Field(f)
		}
		b.fieldsOnly = true
	}
	return b
}

func (b *selectorBuilder) Group(qb QueryBuilder) SelectorBuilder {
	b.parts = append(b.parts, queryPart{partType: selectorGroupPart, sql: "(" + qb.Build().Sql + ")"})
	return b
}

func (b *selectorBuilder) Field(qb QueryBuilder) SelectorBuilder {
	b.parts = append(
		b.parts, queryPart{
			partType: selectorFieldPart,
			sql:      qb.Build().Sql,
			builder:  qb,
		},
	)
	return b
}

func (b *selectorBuilder) As(value string) SelectorBuilder {
	b.parts = append(b.parts, queryPart{partType: selectorAliasPart, sql: fmt.Sprintf("AS %s", value)})
	return b
}

func (b *selectorBuilder) Min(qb QueryBuilder) SelectorBuilder {
	b.parts = append(
		b.parts,
		queryPart{
			partType: selectorAggregatePart,
			sql:      fmt.Sprintf("MIN(%s)", qb.Build().Sql),
			builder:  qb,
		},
	)
	return b
}

func (b *selectorBuilder) Max(qb QueryBuilder) SelectorBuilder {
	b.parts = append(
		b.parts,
		queryPart{
			partType: selectorAggregatePart,
			sql:      fmt.Sprintf("MAX(%s)", qb.Build().Sql),
			builder:  qb,
		},
	)
	return b
}

func (b *selectorBuilder) Avg(qb QueryBuilder) SelectorBuilder {
	b.parts = append(
		b.parts,
		queryPart{
			partType: selectorAggregatePart,
			sql:      fmt.Sprintf("AVG(%s)", qb.Build().Sql),
			builder:  qb,
		},
	)
	return b
}

func (b *selectorBuilder) Count(qb QueryBuilder) SelectorBuilder {
	b.parts = append(
		b.parts,
		queryPart{
			partType: selectorAggregatePart,
			sql:      fmt.Sprintf("COUNT(%s)", qb.Build().Sql),
			builder:  qb,
		},
	)
	return b
}

func (b *selectorBuilder) Sum(qb QueryBuilder) SelectorBuilder {
	b.parts = append(
		b.parts,
		queryPart{
			partType: selectorAggregatePart,
			sql:      fmt.Sprintf("SUM(%s)", qb.Build().Sql),
			builder:  qb,
		},
	)
	return b
}

func (b *selectorBuilder) Build() BuildResult {
	values := make(map[string]any)
	n := len(b.parts)
	sql := make([]string, n)
	for i, part := range b.parts {
		sql[i] = part.sql
	}
	separator := " "
	if b.fieldsOnly {
		separator = ","
	}
	return BuildResult{strings.Join(sql, separator), values}
}
