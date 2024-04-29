package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestRemoveRepository(t *testing.T) {
	te := Entity[testEntity]()
	t.Run(
		"equal id", func(t *testing.T) {
			id := 1
			r := Repository[testEntity, int](nil).Remove(
				Filter().Field(te.Id()).Equal().Value(id, "id"),
			)
			build := r.Build()
			assert.Equal(t, id, build.Values["id"])
			assert.Equal(
				t,
				`DELETE FROM test AS t WHERE t.id = @id RETURNING *`,
				build.Sql,
			)
		},
	)
	t.Run(
		"in ids slice", func(t *testing.T) {
			ids := []int{1, 2, 3}
			r := Repository[testEntity, int](nil).Remove(
				Filter().Field(te.Id()).In().Value(ids, "ids"),
				Selector(te.Id()),
			)
			build := r.Build()
			assert.Equal(t, ids, build.Values["ids"])
			assert.Equal(
				t,
				`DELETE FROM test AS t WHERE t.id IN (@ids) RETURNING t.id`,
				build.Sql,
			)
		},
	)
}
