package peak

type testModel struct {
	Id    int    `db:"id"`
	Email string `db:"email"`
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
