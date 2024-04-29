package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestShapeBuilder(t *testing.T) {
	te := Entity[testEntity]()
	t.Run(
		"start", func(t *testing.T) {
			assert.Equal(
				t,
				"OFFSET 20",
				Shape().Start(20).Build().Sql,
			)
		},
	)
	t.Run(
		"max", func(t *testing.T) {
			assert.Equal(
				t,
				"LIMIT 20",
				Shape().Max(20).Build().Sql,
			)
		},
	)
	t.Run(
		"max and start", func(t *testing.T) {
			assert.Equal(
				t,
				"LIMIT 20 OFFSET 40",
				Shape().Max(20).Start(40).Build().Sql,
			)
		},
	)
	t.Run(
		"manual sort", func(t *testing.T) {
			assert.Equal(
				t,
				"ORDER BY t.id ASC,t.email DESC",
				Shape().Sort().Up(te.Id()).Down(te.Email()).Build().Sql,
			)
		},
	)
	t.Run(
		"auto sort", func(t *testing.T) {
			assert.Equal(
				t,
				"ORDER BY t.id ASC,t.email DESC",
				Shape().Sort().Use(
					Sorter{Field: "t.id", Direction: SortUp}, Sorter{Field: "t.email", Direction: SortDown},
				).Build().Sql,
			)
		},
	)
	t.Run(
		"prevent duplicates", func(t *testing.T) {
			assert.Equal(
				t,
				"GROUP BY t.id,t.email",
				Shape().Duplicates().Prevent(te.Id(), te.Email()).Build().Sql,
			)
		},
	)
	t.Run(
		"remove duplicates", func(t *testing.T) {
			r := Repository[testEntity, []int](nil).
				Find(
					Shape().Duplicates().Remove(),
				)
			assert.Equal(
				t,
				"SELECT DISTINCT t.id,t.email FROM test AS t",
				r.Build().Sql,
			)
		},
	)
	t.Run(
		"aggregate func added to group by", func(t *testing.T) {
			r := Repository[testEntity, []int](nil).
				Find(
					Selector(te.Email()),
					Selector().Count(te.Id()),
				)
			assert.Equal(
				t,
				"SELECT t.email,COUNT(t.id) FROM test AS t GROUP BY t.email",
				r.Build().Sql,
			)
		},
	)
}
