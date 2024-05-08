package peak

import (
	"fmt"
	"reflect"
	"strings"
	
	"github.com/creamsensation/quirk"
)

type FilterBuilder interface {
	QueryBuilder
	Group(builders ...FilterBuilder) FilterBuilder
	Field(field QueryBuilder) FilterBuilder
	Equal(equal ...bool) FilterBuilder
	Is() FilterBuilder
	In(in ...bool) FilterBuilder
	Value(value any, name ...string) FilterBuilder
	Gt() FilterBuilder
	Gte() FilterBuilder
	Lt() FilterBuilder
	Lte() FilterBuilder
	Match() FilterBuilder
	Not(not ...bool) FilterBuilder
	TsQuery(values ...any) FilterBuilder
}

type filterBuilder struct {
	parts []queryPart
	or    bool
	after bool
}

const (
	filterGroupPart    = "filter-group"
	filterFieldPart    = "filter-field"
	filterOperatorPart = "filter-operator"
	filterValuePart    = "filter-value"
)

var (
	fulltextValuesReplacer = strings.NewReplacer(
	
	)
)

func Filter(modifiers ...Modifier) FilterBuilder {
	fb := &filterBuilder{
		parts: make([]queryPart, 0),
	}
	for _, m := range modifiers {
		if m == nil {
			continue
		}
		switch m.Type() {
		case modifierOr:
			fb.or = true
		case modifierAfter:
			fb.after = true
		}
	}
	return fb
}

func (b *filterBuilder) Group(builders ...FilterBuilder) FilterBuilder {
	var sql []string
	values := make(map[string]any)
	for i, item := range builders {
		if item == nil {
			continue
		}
		build := item.Build()
		if i > 0 {
			f := item.(*filterBuilder)
			if !f.or {
				sql = append(sql, "AND")
			}
			if f.or {
				sql = append(sql, "OR")
			}
		}
		sql = append(sql, build.Sql)
		for k, v := range build.Values {
			values[k] = v
		}
	}
	b.parts = append(
		b.parts,
		queryPart{
			partType: filterGroupPart,
			sql:      "(" + strings.Join(sql, " ") + ")",
			value:    values,
		},
	)
	return b
}

func (b *filterBuilder) Field(field QueryBuilder) FilterBuilder {
	b.parts = append(
		b.parts,
		queryPart{
			partType: filterFieldPart,
			sql:      field.Build().Sql,
		},
	)
	return b
}

func (b *filterBuilder) Equal(equal ...bool) FilterBuilder {
	isEqual := true
	operator := "="
	if len(equal) > 0 {
		isEqual = equal[0]
	}
	if !isEqual {
		operator = "!="
	}
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: operator})
	return b
}

func (b *filterBuilder) Match() FilterBuilder {
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: "@@"})
	return b
}

func (b *filterBuilder) TsQuery(values ...any) FilterBuilder {
	queryName := "query" + generateRandomString(8)
	b.parts = append(
		b.parts, queryPart{
			partType: filterOperatorPart,
			sql:      "to_tsquery(@" + queryName + ")",
			name:     queryName,
			value:    quirk.CreateTsQuery(values...),
		},
	)
	return b
}

func (b *filterBuilder) Not(not ...bool) FilterBuilder {
	if len(not) > 0 && !not[0] {
		return b
	}
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: "NOT"})
	return b
}

func (b *filterBuilder) Is() FilterBuilder {
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: "IS"})
	return b
}

func (b *filterBuilder) In(in ...bool) FilterBuilder {
	isIn := true
	operator := "IN"
	if len(in) > 0 {
		isIn = in[0]
	}
	if !isIn {
		operator = "NOT IN"
	}
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: operator})
	return b
}

func (b *filterBuilder) Value(value any, name ...string) FilterBuilder {
	valueName := generateRandomString(8)
	if len(name) > 0 {
		valueName = name[0]
	}
	sql := fmt.Sprintf("@%s", valueName)
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		sql = "(" + sql + ")"
	}
	b.parts = append(
		b.parts,
		queryPart{
			partType: filterValuePart,
			sql:      sql,
			name:     valueName,
			value:    value,
		},
	)
	return b
}

func (b *filterBuilder) Gt() FilterBuilder {
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: ">"})
	return b
}

func (b *filterBuilder) Gte() FilterBuilder {
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: ">="})
	return b
}

func (b *filterBuilder) Lt() FilterBuilder {
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: "<"})
	return b
}

func (b *filterBuilder) Lte() FilterBuilder {
	b.parts = append(b.parts, queryPart{partType: filterOperatorPart, sql: "<="})
	return b
}

func (b *filterBuilder) Build() BuildResult {
	values := make(map[string]any)
	n := len(b.parts)
	sql := make([]string, n)
	for i, part := range b.parts {
		sql[i] = part.sql
		if part.value != nil {
			switch v := part.value.(type) {
			case map[string]any:
				for key, value := range v {
					values[key] = value
				}
			default:
				values[part.name] = part.value
			}
		}
	}
	return BuildResult{strings.TrimSpace(strings.Join(sql, " ")), values}
}
