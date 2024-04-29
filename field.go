package peak

import (
	"fmt"
)

type Field interface {
	QueryBuilder
	Default(value any) Field
	NotNull(notNull ...bool) Field
	Prefix(prefix string) Field
	PrimaryKey(primaryKey ...bool) Field
	Type(fieldType string) Field
	Unique(unique ...bool) Field
	Relationship(relationship Field) Field
}

type field struct {
	Field
	defaultValue string
	fieldType    string
	name         string
	notNull      bool
	table        string
	prefix       string
	primaryKey   bool
	relationship *field
	unique       bool
}

func (f *field) Default(value any) Field {
	switch v := value.(type) {
	case string:
		f.defaultValue = fmt.Sprintf("'%s'", v)
	default:
		f.defaultValue = fmt.Sprintf("%v", value)
	}
	return f
}

func (f *field) NotNull(notNull ...bool) Field {
	n := len(notNull)
	if n == 0 {
		f.notNull = true
	}
	if n > 0 {
		f.notNull = notNull[0]
	}
	return f
}

func (f *field) Prefix(prefix string) Field {
	f.prefix = prefix
	return f
}

func (f *field) PrimaryKey(primaryKey ...bool) Field {
	n := len(primaryKey)
	if n == 0 {
		f.primaryKey = true
	}
	if n > 0 {
		f.primaryKey = primaryKey[0]
	}
	f.notNull = true
	return f
}

func (f *field) Relationship(relationship Field) Field {
	f.relationship = relationship.(*field)
	return f
}

func (f *field) Type(fieldType string) Field {
	f.fieldType = fieldType
	return f
}

func (f *field) Unique(unique ...bool) Field {
	n := len(unique)
	if n == 0 {
		f.unique = true
	}
	if n > 0 {
		f.unique = unique[0]
	}
	return f
}

func (f *field) Build() BuildResult {
	return BuildResult{f.prefix + "." + f.name, nil}
}
