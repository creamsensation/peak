package peak

import "reflect"

type valuesBuilder struct {
	values any
}

const (
	fieldDbTagName = "db"
)

func Use(values any) QueryBuilder {
	return &valuesBuilder{
		values: values,
	}
}

func (b *valuesBuilder) Build() BuildResult {
	values := make(map[string]any)
	t := reflect.TypeOf(b.values)
	v := reflect.ValueOf(b.values)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if t.Kind() == reflect.Map {
		mapValues, ok := b.values.(map[string]any)
		if !ok {
			panic(ErrorInvalidMap)
		}
		return BuildResult{Values: mapValues}
	}
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		kind := t.Field(i).Type.Kind()
		dbName := t.Field(i).Tag.Get(fieldDbTagName)
		if len(dbName) == 0 {
			continue
		}
		f := v.FieldByName(name)
		if !f.IsValid() || (f.IsValid() && f.IsZero() && kind != reflect.Bool) {
			values[dbName] = nil
			continue
		}
		values[dbName] = f.Interface()
	}
	return BuildResult{Values: values}
}
