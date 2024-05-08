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
	t.Run(
		"insert with vectors", func(t *testing.T) {
			r := Repository[fulltextEntity, int](nil).Save(
				Use(
					fulltextModel{
						Name:    "Dominik",
						Vectors: TsVector("Dominik"),
					},
				),
				Selector(te.Id()),
			)
			b := r.Build()
			assert.Equal(
				t,
				`INSERT INTO fulltexts AS f (f.name,f.vectors) VALUES (@name,to_tsvector(@vectors)) RETURNING t.id`,
				b.Sql,
			)
			assert.Equal(t, "dominik", b.Values["vectors"])
		},
	)
	t.Run(
		"update with vectors", func(t *testing.T) {
			r := Repository[fulltextEntity, int](nil).Save(
				Use(
					fulltextModel{
						Id:   1,
						Name: "Dominik",
					},
				),
				Selector(te.Id()),
			)
			b := r.Build()
			assert.Equal(
				t,
				`UPDATE fulltexts AS f SET f.name = @name,f.vectors = to_tsvector(@vectors) WHERE f.id = @id RETURNING t.id`,
				b.Sql,
			)
			assert.Equal(t, "dominik", b.Values["vectors"])
		},
	)
	t.Run(
		"insert with timestamp", func(t *testing.T) {
			r := Repository[timeEntity, int](nil).Save(
				Selector(te.Id()),
			)
			b := r.Build()
			assert.Equal(
				t,
				`INSERT INTO times AS t (t.created_at,t.updated_at) VALUES (CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) RETURNING t.id`,
				b.Sql,
			)
		},
	)
	t.Run(
		"update with timestamp", func(t *testing.T) {
			r := Repository[timeEntity, int](nil).Save(
				Use(timeModel{Id: 1}),
				Selector(te.Id()),
			)
			b := r.Build()
			assert.Equal(
				t,
				`UPDATE times AS t SET t.updated_at = CURRENT_TIMESTAMP WHERE t.id = @id RETURNING t.id`,
				b.Sql,
			)
		},
	)
}
