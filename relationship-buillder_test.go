package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestRelationshipBuilder(t *testing.T) {
	// te := Entity[testEntity]()
	t.Run(
		"default join", func(t *testing.T) {
			r := Repository[chapterEntity, []chapterModel](nil).Find()
			assert.Equal(
				t,
				`SELECT ch.id,ch.book_id FROM chapters AS ch LEFT JOIN books AS b ON b.id = ch.book_id`,
				r.Build().Sql,
			)
		},
	)
	t.Run(
		"override join", func(t *testing.T) {
			r := Repository[chapterEntity, []chapterModel](nil).Find(
				Relationship(be.Id()).Intersect(),
			)
			assert.Equal(
				t,
				`SELECT ch.id,ch.book_id FROM chapters AS ch INNER JOIN books AS b ON b.id = ch.book_id`,
				r.Build().Sql,
			)
		},
	)
}
