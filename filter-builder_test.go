package peak

import (
	"fmt"
	"strings"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestFilterBuilder(t *testing.T) {
	te := Entity[testEntity]()
	t.Run(
		"in", func(t *testing.T) {
			value := []int{1, 2, 3}
			condition := Filter().
				Field(te.Id()).
				In().
				Value(value)
			r := condition.Build()
			for k, v := range r.Values {
				assert.Equal(
					t,
					fmt.Sprintf("%s.id IN (@%s)", te.Alias(), k),
					r.Sql,
				)
				assert.Equal(t, value, v)
			}
		},
	)
	t.Run(
		"group", func(t *testing.T) {
			value := 1
			condition := Filter().Group(
				Filter().Field(te.Id()).Equal().Value(value),
			)
			r := condition.Build()
			for k, v := range r.Values {
				assert.Equal(
					t,
					fmt.Sprintf("(%s.id = @%s)", te.Alias(), k),
					r.Sql,
				)
				assert.Equal(t, value, v)
			}
		},
	)
	t.Run(
		"or", func(t *testing.T) {
			value1 := 1
			value2 := 2
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Equal().Value(value1, "value1"),
				Filter(Or()).Field(te.Id()).Equal().Value(value2, "value2"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id = @value1 OR t.id = @value2`,
				b.Sql,
			)
		},
	)
	t.Run(
		"or", func(t *testing.T) {
			value1 := 1
			value2 := 2
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Equal().Value(value1, "value1"),
				Filter(Or()).Field(te.Id()).Equal().Value(value2, "value2"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id = @value1 OR t.id = @value2`,
				b.Sql,
			)
		},
	)
	t.Run(
		"gt", func(t *testing.T) {
			value := 1
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Gt().Value(value, "value"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id > @value`,
				b.Sql,
			)
		},
	)
	t.Run(
		"gte", func(t *testing.T) {
			value := 1
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Gte().Value(value, "value"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id >= @value`,
				b.Sql,
			)
		},
	)
	t.Run(
		"lt", func(t *testing.T) {
			value := 1
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Lt().Value(value, "value"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id < @value`,
				b.Sql,
			)
		},
	)
	t.Run(
		"lte", func(t *testing.T) {
			value := 1
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Lte().Value(value, "value"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id <= @value`,
				b.Sql,
			)
		},
	)
	t.Run(
		"aggregation", func(t *testing.T) {
			value := 1
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Id()).Equal().Value(value, "value1"),
				Filter(After()).Field(Selector().Count(te.Id())).Gt().Value(value, "value2"),
			)
			b := r.Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id = @value1 HAVING COUNT(t.id) > @value2`,
				b.Sql,
			)
		},
	)
	t.Run(
		"tsquery", func(t *testing.T) {
			r := Repository[testEntity, []int](nil).Find(
				Filter().Field(te.Vectors()).Match().TsQuery("test1", "test2"),
			)
			b := r.Build()
			qp := `SELECT t.id,t.email FROM test AS t WHERE t.vectors @@ to_tsquery(@query`
			paramName := strings.TrimSuffix(
				strings.TrimPrefix(b.Sql, qp),
				")",
			)
			assert.Contains(
				t,
				b.Sql,
				qp,
			)
			assert.Equal(t, "test:*", b.Values["query"+paramName])
		},
	)
}
