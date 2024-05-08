package peak

import (
	"fmt"
	"strings"
	
	"github.com/thanhpk/randstr"
	
	"github.com/creamsensation/quirk"
)

func TsVector(values ...any) string {
	return quirk.CreateTsVector(values...)
}

func generateRandomString(length int) string {
	// b := make([]byte, length)
	// _, err := rand.Read(b)
	// if err != nil {
	// 	return ""
	// }
	return randstr.Hex(length)
}

func findField(name string, fields []Field) *field {
	for _, item := range fields {
		if item == nil {
			continue
		}
		f := item.(*field)
		if f.name == name {
			return f
		}
	}
	return nil
}

func createInsertSqlFromValues(fields []Field, values map[string]any) string {
	sql := make([]string, 0)
	for _, item := range fields {
		if item == nil {
			continue
		}
		f := item.(*field)
		if f.primaryKey {
			continue
		}
		v, ok := values[f.name]
		if !ok {
			continue
		}
		switch val := v.(type) {
		case Safe:
			sql = append(sql, string(val))
		case nil:
			if f.notNull {
				sql = append(sql, "DEFAULT")
				continue
			}
			if !f.notNull {
				sql = append(sql, "NULL")
				continue
			}
		default:
			sql = append(sql, createFieldValuePlaceholder(f.name, f))
		}
	}
	return strings.Join(sql, ",")
}

func createUpdateSqlFromValues(fields []Field, values map[string]any) string {
	sql := make([]string, 0)
	for _, item := range fields {
		if item == nil {
			continue
		}
		f := item.(*field)
		if f.primaryKey {
			continue
		}
		v, ok := values[f.name]
		if !ok {
			continue
		}
		switch val := v.(type) {
		case Safe:
			sql = append(sql, fmt.Sprintf("%s.%s = %s", f.prefix, f.name, string(val)))
		case nil:
			if !f.notNull {
				sql = append(sql, fmt.Sprintf("%s.%s = NULL", f.prefix, f.name))
				continue
			}
		default:
			sql = append(sql, fmt.Sprintf("%s.%s = %s", f.prefix, f.name, createFieldValuePlaceholder(f.name, f)))
		}
	}
	return strings.Join(sql, ",")
}

func createFieldValuePlaceholder(k string, f *field) string {
	switch f.dataType {
	case TsVectorDataType:
		return fmt.Sprintf("to_tsvector(@%s)", k)
	default:
		return fmt.Sprintf("@%s", k)
	}
}

func buildFieldsSql[T QueryBuilder](builders ...T) string {
	var r []string
	for _, b := range builders {
		r = append(r, b.Build().Sql)
	}
	return strings.Join(r, ",")
}

func buildFieldsSqlWithoutPrimaryKey[T QueryBuilder](builders ...T) string {
	var r []string
	for _, b := range builders {
		switch v := any(b).(type) {
		case *field:
			if !v.primaryKey {
				r = append(r, b.Build().Sql)
			}
		default:
			r = append(r, b.Build().Sql)
		}
	}
	return strings.Join(r, ",")
}

func buildJoins(q *sqlBuilder, relationships []*relationshipBuilder, fields []Field) {
	for _, item := range fields {
		f := item.(*field)
		if f.relationship == nil {
			continue
		}
		var rb *relationshipBuilder
		for _, r := range relationships {
			if r.field.prefix == f.relationship.prefix && r.field.name == f.relationship.name {
				rb = r
			}
		}
		if rb != nil {
			q.Q(
				buildJoinSql(rb.joinType, f).Sql,
			)
			continue
		}
		q.Q(
			buildJoinSql(leftJoin, f).Sql,
		)
	}
}

func buildBeforeAggregationFilters(
	q *sqlBuilder, filters []*filterBuilder, values *map[string]any,
) {
	beforeAggregationFilters := make([]*filterBuilder, 0)
	for _, f := range filters {
		if f.after {
			continue
		}
		beforeAggregationFilters = append(beforeAggregationFilters, f)
	}
	buildFilters(q, beforeAggregationFilters, values, "WHERE")
}
func buildAfterAggregationFilters(
	q *sqlBuilder, filters []*filterBuilder, values *map[string]any,
) {
	afterAggregationFilters := make([]*filterBuilder, 0)
	for _, f := range filters {
		if !f.after {
			continue
		}
		afterAggregationFilters = append(afterAggregationFilters, f)
	}
	buildFilters(q, afterAggregationFilters, values, "HAVING")
}

func buildFilters(
	q *sqlBuilder, filters []*filterBuilder, values *map[string]any, keyword string,
) {
	if len(filters) == 0 {
		return
	}
	for i, f := range filters {
		q = q.If(i == 0, keyword)
		if i > 0 {
			q = q.If(!f.or, "AND")
			q = q.If(f.or, "OR")
		}
		filterResult := f.Build()
		q = q.Q(filterResult.Sql)
		for k, v := range filterResult.Values {
			(*values)[k] = v
		}
	}
}

func buildGroupShapes(shapes []*shapeBuilder) string {
	for _, s := range shapes {
		if len(s.groupFields) == 0 {
			continue
		}
		return s.Build().Sql
	}
	return ""
}

func buildNonGroupShapes(shapes []*shapeBuilder) string {
	for _, s := range shapes {
		if len(s.groupFields) > 0 {
			continue
		}
		return s.Build().Sql
	}
	return ""
}

func doesExistDistinct(shapes []*shapeBuilder) bool {
	for _, s := range shapes {
		if s.distinct {
			return true
		}
	}
	return false
}

func getPrimaryKeyField(fields ...Field) *field {
	for _, item := range fields {
		f := item.(*field)
		if f.primaryKey {
			return f
		}
	}
	return nil
}
