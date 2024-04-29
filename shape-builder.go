package peak

import (
	"fmt"
	"strings"
)

type ShapeBuilder interface {
	QueryBuilder
	Start(start int) ShapeBuilder
	Max(max int) ShapeBuilder
	Duplicates() DuplicatesBuilder
	Sort(aliases ...map[string]string) SortBuilder
}

type SortBuilder interface {
	QueryBuilder
	Use(sorters ...Sorter) SortBuilder
	Up(field QueryBuilder) SortBuilder
	Down(field QueryBuilder) SortBuilder
}

type DuplicatesBuilder interface {
	QueryBuilder
	Prevent(fields ...QueryBuilder) ShapeBuilder
	Remove(remove ...bool) ShapeBuilder
}

type shapeBuilder struct {
	start       int
	max         int
	groupFields []QueryBuilder
	distinct    bool
	sorts       []Sorter
	sortAliases map[string]string
}

type Sorter struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

const (
	SortUp   = "ASC"
	SortDown = "DESC"
)

func Shape() ShapeBuilder {
	return &shapeBuilder{
		sorts:       make([]Sorter, 0),
		sortAliases: make(map[string]string),
	}
}

func (b *shapeBuilder) Start(start int) ShapeBuilder {
	b.start = start
	return b
}

func (b *shapeBuilder) Max(max int) ShapeBuilder {
	b.max = max
	return b
}

func (b *shapeBuilder) Sort(aliases ...map[string]string) SortBuilder {
	if len(aliases) > 0 {
		b.sortAliases = aliases[0]
	}
	return b
}

func (b *shapeBuilder) Up(field QueryBuilder) SortBuilder {
	b.sorts = append(b.sorts, Sorter{Field: field.Build().Sql, Direction: SortUp})
	return b
}

func (b *shapeBuilder) Down(field QueryBuilder) SortBuilder {
	b.sorts = append(b.sorts, Sorter{Field: field.Build().Sql, Direction: SortDown})
	return b
}

func (b *shapeBuilder) Use(sorters ...Sorter) SortBuilder {
	b.sorts = append(b.sorts, sorters...)
	return b
}

func (b *shapeBuilder) Duplicates() DuplicatesBuilder {
	return b
}

func (b *shapeBuilder) Prevent(fields ...QueryBuilder) ShapeBuilder {
	b.groupFields = append(b.groupFields, fields...)
	return b
}

func (b *shapeBuilder) Remove(distinct ...bool) ShapeBuilder {
	b.distinct = true
	if len(distinct) > 0 {
		b.distinct = distinct[0]
	}
	return b
}

func (b *shapeBuilder) Build() BuildResult {
	q := createSqlBuilder().
		If(len(b.groupFields) > 0, fmt.Sprintf("GROUP BY %s", buildFieldsSql(b.groupFields...))).
		If(len(b.sorts) > 0, fmt.Sprintf("ORDER BY %s", b.buildSorts())).
		If(b.max > 0, fmt.Sprintf("LIMIT %d", b.max)).
		If(b.start > 0, fmt.Sprintf("OFFSET %d", b.start))
	return BuildResult{q.Build(), nil}
}

func (b *shapeBuilder) buildSorts() string {
	r := make([]string, len(b.sorts))
	for i, s := range b.sorts {
		f, ok := b.sortAliases[s.Field]
		if ok {
			r[i] = fmt.Sprintf("%s %s", f, s.Direction)
			continue
		}
		r[i] = fmt.Sprintf("%s %s", s.Field, s.Direction)
	}
	return strings.Join(r, ",")
}
