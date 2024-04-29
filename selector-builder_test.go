package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestSelector(t *testing.T) {
	te := Entity[testEntity]()
	t.Run(
		"as", func(t *testing.T) {
			s := Selector().
				Field(te.Id()).
				As("test_id")
			assert.Equal(
				t,
				"t.id AS test_id",
				s.Build().Sql,
			)
		},
	)
	t.Run(
		"fields only", func(t *testing.T) {
			s := Selector(
				te.Email(),
				te.Id(),
			)
			assert.Equal(t, "t.email,t.id", s.Build().Sql)
		},
	)
}
