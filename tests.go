package peak

import (
	"time"
	
	"github.com/creamsensation/quirk"
)

type testModel struct {
	Id    int    `db:"id"`
	Email string `db:"email"`
}

type fulltextModel struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Vectors string `db:"vectors"`
}

type timeModel struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

var (
	te  = Entity[testEntity]()
	be  = Entity[bookEntity]()
	che = Entity[chapterEntity]()
)

// test entity

type testEntity struct {
	EntityBuilder
}

func (e testEntity) Table() string {
	return "test"
}

func (e testEntity) Alias() string {
	return "t"
}

func (e testEntity) Fields() []Field {
	return e.Register(
		e.Id(),
		e.Email(),
	)
}

func (e testEntity) Id() Field {
	return e.Field("id").
		Type("SERIAL").
		PrimaryKey()
}

func (e testEntity) Email() Field {
	return e.Field("email").
		Type("VARCHAR(255)").
		NotNull()
}

func (e testEntity) Vectors() Field {
	return e.Field("vectors").
		Type("TSVECTOR").
		NotNull().
		Default("")
}

type bookModel struct {
	Id       int            `db:"id"`
	Chapters []chapterModel `db:"chapters"`
}

// test book entity

type bookEntity struct {
	EntityBuilder
}

func (e bookEntity) Table() string {
	return "books"
}

func (e bookEntity) Alias() string {
	return "b"
}

func (e bookEntity) Fields() []Field {
	return e.Register(
		e.Id(),
	)
}

func (e bookEntity) Id() Field {
	return e.Field("id").
		Type("SERIAL").
		PrimaryKey()
}

type chapterModel struct {
	Id     int `db:"id"`
	BookId int `db:"book_id"`
}

// test chapter entity

type chapterEntity struct {
	EntityBuilder
}

func (e chapterEntity) Table() string {
	return "chapters"
}

func (e chapterEntity) Alias() string {
	return "ch"
}

func (e chapterEntity) Fields() []Field {
	return e.Register(
		e.Id(),
		e.BookId(),
	)
}

func (e chapterEntity) Id() Field {
	return e.Field("id").
		Type("SERIAL").
		PrimaryKey()
}

func (e chapterEntity) BookId() Field {
	return e.Field("book_id").
		Type("INT").
		Relationship(be.Id())
}

// test fulltext entity

type fulltextEntity struct {
	EntityBuilder
}

func (e fulltextEntity) Table() string {
	return "fulltexts"
}

func (e fulltextEntity) Alias() string {
	return "f"
}

func (e fulltextEntity) Fields() []Field {
	return e.Register(
		e.Id(),
		e.Name(),
		e.Vectors(),
	)
}

func (e fulltextEntity) Id() Field {
	return e.Field("id").
		Type("SERIAL").
		PrimaryKey()
}

func (e fulltextEntity) Name() Field {
	return e.Field("name").
		Type("VARCHAR(255)").
		NotNull()
}

func (e fulltextEntity) Vectors() Field {
	return e.Field("vectors").
		TsVector().
		NotNull().
		Default("").
		ValueFactory(
			func(operation string, values Map) Value {
				if operation == Insert {
					return nil
				}
				return TsVector(
					values[e.Name().Name()],
				)
			},
		)
}

// test time entity

type timeEntity struct {
	EntityBuilder
}

func (e timeEntity) Table() string {
	return "times"
}

func (e timeEntity) Alias() string {
	return "t"
}

func (e timeEntity) Fields() []Field {
	return e.Register(
		e.Id(),
		e.CreatedAt(),
		e.UpdatedAt(),
	)
}

func (e timeEntity) Id() Field {
	return e.Field("id").
		Type("SERIAL").
		PrimaryKey()
}

func (e timeEntity) CreatedAt() Field {
	return e.Field("created_at").
		Type("TIMESTAMP").
		NotNull().
		Default(quirk.CurrentTimestamp).
		ValueFactory(
			func(operation string, values Map) Value {
				if operation == Update {
					return nil
				}
				return Safe(quirk.CurrentTimestamp)
			},
		)
}

func (e timeEntity) UpdatedAt() Field {
	return e.Field("updated_at").
		Type("TIMESTAMP").
		NotNull().
		Default(quirk.CurrentTimestamp).
		ValueFactory(
			func(operation string, values Map) Value {
				return Safe(quirk.CurrentTimestamp)
			},
		)
}
