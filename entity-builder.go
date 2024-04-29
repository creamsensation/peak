package peak

import "reflect"

type EntityBuilder struct {
	table  string
	prefix string
}

var (
	entityBuilderFieldName = reflect.TypeOf(EntityBuilder{}).Name()
)

func (b EntityBuilder) Field(name string) Field {
	return &field{table: b.table, prefix: b.prefix, name: name}
}

func (b EntityBuilder) Register(field ...Field) []Field {
	return field
}
