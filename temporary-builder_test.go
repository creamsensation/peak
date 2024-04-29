package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestTemporaryBuilder(t *testing.T) {
	r := Repository[chapterEntity, int](nil).Save(
		Temporary(
			"books",
			Repository[bookEntity, int](nil).Find(),
		),
		Use(chapterModel{Id: 1}),
		Use(Map{"book_id": Safe("books.id")}),
		Selector(che.Id()),
	)
	assert.Equal(
		t,
		`WITH (SELECT b.id FROM books AS b) AS books UPDATE chapters AS ch SET ch.book_id = books.id WHERE ch.id = @id RETURNING ch.id`,
		r.Build().Sql,
	)
}
