package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestFindRepository(t *testing.T) {
	te := Entity[testEntity]()
	t.Run(
		"basic", func(t *testing.T) {
			r := Repository[testEntity, []int](nil).Find().Build()
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t`,
				r.Sql,
			)
		},
	)
	t.Run(
		"with filters", func(t *testing.T) {
			v1 := 1
			v2 := 2
			r := Repository[testEntity, []int](nil).
				Find(
					Filter().Field(te.Id()).Equal().Value(v1, "v1"),
					Filter(Or()).Group(
						Filter().Field(te.Id()).Equal().Value(v2, "v2"),
					),
				).
				Build()
			assert.Equal(t, v1, r.Values["v1"])
			assert.Equal(t, v2, r.Values["v2"])
			assert.Equal(
				t,
				`SELECT t.id,t.email FROM test AS t WHERE t.id = @v1 OR (t.id = @v2)`,
				r.Sql,
			)
		},
	)
}
