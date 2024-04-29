package peak

import (
	"reflect"
	"sync"
	
	"github.com/creamsensation/quirk"
)

type Runner interface {
	Add(builders ...QueryBuilder)
	AllInSequence(target ...any) error
	AllInParallel(target ...any) error
	MustAllInSequence(target ...any)
	MustAllInParallel(target ...any)
}

type runner struct {
	db       *quirk.DB
	builders []QueryBuilder
	wg       *sync.WaitGroup
}

func Run(db *quirk.DB) Runner {
	return &runner{
		db:       db,
		builders: make([]QueryBuilder, 0),
	}
}

func (r *runner) Add(builders ...QueryBuilder) {
	r.builders = append(r.builders, builders...)
}

func (r *runner) AllInSequence(target ...any) error {
	if len(r.builders) != len(target) {
		return ErrorMismatchQueryTarget
	}
	if r.db == nil {
		return ErrorMissingDatabase
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	for i, item := range r.builders {
		b := item.Build()
		if reflect.TypeOf(target[i]).Kind() != reflect.Ptr {
			return ErrorTargetNoPtr
		}
		if err := tx.Q(b.Sql, b.Values).Exec(target[i]); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	return tx.Commit()
}

func (r *runner) AllInParallel(target ...any) error {
	if len(r.builders) != len(target) {
		return ErrorMismatchQueryTarget
	}
	if r.db == nil {
		return ErrorMissingDatabase
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	r.wg = new(sync.WaitGroup)
	errs := make(chan error, len(r.builders))
	for i, item := range r.builders {
		r.wg.Add(1)
		go func(i int, item QueryBuilder) {
			defer r.wg.Done()
			b := item.Build()
			if reflect.TypeOf(target[i]).Kind() != reflect.Ptr {
				errs <- ErrorTargetNoPtr
				return
			}
			if err := r.db.Q(b.Sql, b.Values).Exec(target[i]); err != nil {
				errs <- err
				return
			}
		}(i, item)
	}
	r.wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				return rerr
			}
			return err
		}
	}
	return tx.Commit()
}

func (r *runner) MustAllInSequence(target ...any) {
	err := r.AllInSequence(target...)
	if err != nil {
		panic(err)
	}
}

func (r *runner) MustAllInParallel(target ...any) {
	err := r.AllInParallel(target...)
	if err != nil {
		panic(err)
	}
}
