package peak

import (
	"reflect"
)

type entity interface {
	Table() string
	Alias() string
	Fields() []Field
}

func Entity[E entity]() *E {
	e := new(E)
	p := any(e).(entity)
	v := reflect.ValueOf(e)
	f := v.Elem().FieldByName(entityBuilderFieldName)
	if !f.IsValid() {
		return e
	}
	f.Set(reflect.ValueOf(EntityBuilder{table: p.Table(), prefix: p.Alias()}))
	return e
}
