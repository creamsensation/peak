package peak

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestEntityMigrator(t *testing.T) {
	upSql := Migrate[testEntity](nil).GetUpSql()
	downSql := Migrate[testEntity](nil).Cascade().GetDownSql()
	assert.Equal(
		t,
		`CREATE TABLE IF NOT EXISTS test (
	id SERIAL NOT NULL PRIMARY KEY,
	email VARCHAR(255) NOT NULL
)`,
		upSql,
	)
	assert.Equal(t, "DROP TABLE IF EXISTS test CASCADE", downSql)
}
