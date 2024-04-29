package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestSaveRepository(t *testing.T) {
	te := Entity[testEntity]()
	t.Run(
		"insert", func(t *testing.T) {
			r := Repository[testEntity, int](nil).Save(
				Use(testModel{Email: "test@test.com"}),
				Selector(te.Id()),
			)
			assert.Equal(
				t,
				`INSERT INTO test AS t (t.email) VALUES (@email) RETURNING t.id`,
				r.Build().Sql,
			)
		},
	)
	t.Run(
		"update", func(t *testing.T) {
			r := Repository[testEntity, int](nil).Save(
				Use(testModel{Id: 1, Email: "test@test.com"}),
				Selector(te.Id()),
			)
			assert.Equal(
				t,
				`UPDATE test AS t SET t.email = @email WHERE t.id = @id RETURNING t.id`,
				r.Build().Sql,
			)
		},
	)
}
