package peak

import (
	"github.com/creamsensation/quirk"
)

type RepositoryManager[E entity, R result] interface {
	Find(builders ...QueryBuilder) FindRepository[R]
	Save(builders ...QueryBuilder) SaveRepository[R]
	Remove(builders ...QueryBuilder) RemoveRepository[R]
}

type repository[E entity, R result] struct {
	db     *quirk.DB
	entity *E
}

type result interface{}

func Repository[E entity, R result](db *quirk.DB) RepositoryManager[E, R] {
	return &repository[E, R]{
		db:     db,
		entity: Entity[E](),
	}
}

func (r *repository[E, R]) Find(builders ...QueryBuilder) FindRepository[R] {
	tree := createTree(builders...)
	return &findRepository[E, R]{
		repository:    r,
		filters:       tree.filters,
		relationships: tree.relationships,
		selectors:     tree.selectors,
		shapes:        tree.shapes,
	}
}

func (r *repository[E, R]) Save(builders ...QueryBuilder) SaveRepository[R] {
	tree := createTree(builders...)
	return &saveRepository[E, R]{
		repository:    r,
		filters:       tree.filters,
		relationships: tree.relationships,
		selectors:     tree.selectors,
		temporaries:   tree.temporaries,
		values:        tree.values,
	}
}

func (r *repository[E, R]) Remove(builders ...QueryBuilder) RemoveRepository[R] {
	tree := createTree(builders...)
	return &removeRepository[E, R]{
		repository: r,
		filters:    tree.filters,
		selectors:  tree.selectors,
	}
}
