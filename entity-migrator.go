package peak

import (
	"fmt"
	"strings"
	
	"github.com/creamsensation/quirk"
)

type EntityMigrator interface {
	Cascade(cascade ...bool) EntityMigrator
	GetUpSql() string
	GetDownSql() string
	Up() error
	Down() error
	
	MustUp()
	MustDown()
}

type entityMigrator struct {
	db      *quirk.DB
	entity  entity
	cascade bool
}

func Migrate[E entity](db *quirk.DB) EntityMigrator {
	return &entityMigrator{
		db:     db,
		entity: any(new(E)).(entity),
	}
}

func (m *entityMigrator) Cascade(cascade ...bool) EntityMigrator {
	m.cascade = true
	if len(cascade) > 0 {
		m.cascade = cascade[0]
	}
	return m
}

func (m *entityMigrator) GetUpSql() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n)", m.entity.Table(), m.createFieldsSql())
}

func (m *entityMigrator) GetDownSql() string {
	return createSqlBuilder().
		Q(fmt.Sprintf("DROP TABLE IF EXISTS %s", m.entity.Table())).
		If(m.cascade, "CASCADE").
		Build()
}

func (m *entityMigrator) Up() error {
	return nil
}

func (m *entityMigrator) Down() error {
	return nil
}

func (m *entityMigrator) MustUp() {
	err := m.Up()
	if err != nil {
		panic(err)
	}
}

func (m *entityMigrator) MustDown() {
	err := m.Down()
	if err != nil {
		panic(err)
	}
}

func (m *entityMigrator) createFieldsSql() string {
	fields := m.entity.Fields()
	r := make([]string, len(fields))
	for i, item := range fields {
		f := item.(*field)
		sql := createSqlBuilder().
			Q("\t"+f.name).
			Q(f.dataType).
			If(f.notNull, "NOT NULL").
			If(len(f.defaultValue) > 0, "DEFAULT "+f.defaultValue).
			If(f.unique, "UNIQUE").
			If(f.primaryKey, "PRIMARY KEY")
		if f.relationship != nil {
			sql = sql.Q("REFERENCES " + f.relationship.table + "(" + f.relationship.name + ")")
		}
		r[i] = sql.Build()
	}
	return strings.Join(r, ",\n")
}
