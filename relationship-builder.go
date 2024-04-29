package peak

import "fmt"

type RelationshipBuilder interface {
	QueryBuilder
	Keep() RelationshipBuilder
	Prefer() RelationshipBuilder
	Intersect() RelationshipBuilder
	Unify() RelationshipBuilder
	Combine() RelationshipBuilder
}

type relationshipBuilder struct {
	field    *field
	joinType string
}

const (
	leftJoin  = "LEFT"
	rightJoin = "RIGHT"
	innerJoin = "INNER"
	fullJoin  = "FULL"
	crossJoin = "CROSS"
)

func Relationship(relationshipField Field) RelationshipBuilder {
	f := relationshipField.(*field)
	return &relationshipBuilder{
		field:    f,
		joinType: leftJoin,
	}
}

func (b *relationshipBuilder) Keep() RelationshipBuilder {
	b.joinType = leftJoin
	return b
}

func (b *relationshipBuilder) Prefer() RelationshipBuilder {
	b.joinType = rightJoin
	return b
}

func (b *relationshipBuilder) Intersect() RelationshipBuilder {
	b.joinType = innerJoin
	return b
}

func (b *relationshipBuilder) Unify() RelationshipBuilder {
	b.joinType = fullJoin
	return b
}

func (b *relationshipBuilder) Combine() RelationshipBuilder {
	b.joinType = crossJoin
	return b
}

func (b *relationshipBuilder) Build() BuildResult {
	return BuildResult{"", nil}
}

func buildJoinSql(joinType string, field *field) BuildResult {
	sql := fmt.Sprintf(
		"%s JOIN %s AS %s ON %s.%s = %s.%s",
		joinType,
		field.relationship.table,
		field.relationship.prefix,
		field.relationship.prefix,
		field.relationship.name,
		field.prefix,
		field.name,
	)
	return BuildResult{sql, nil}
}
