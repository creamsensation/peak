package peak

import "strings"

type sqlBuilder struct {
	parts []string
}

func createSqlBuilder() *sqlBuilder {
	return &sqlBuilder{
		parts: make([]string, 0),
	}
}

func (q *sqlBuilder) Q(value string) *sqlBuilder {
	q.parts = append(q.parts, value)
	return q
}

func (q *sqlBuilder) If(condition bool, value string) *sqlBuilder {
	if !condition {
		return q
	}
	q.parts = append(q.parts, value)
	return q
}

func (q *sqlBuilder) Build() string {
	return strings.Join(q.parts, " ")
}
