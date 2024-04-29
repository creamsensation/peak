package peak

import "errors"

var (
	ErrorMissingPrimaryKey   = errors.New("missing primary key")
	ErrorMissingDatabase     = errors.New("missing database connection")
	ErrorMismatchQueryTarget = errors.New("mismatch query target len with builders")
	ErrorTargetNoPtr         = errors.New("query target is not a pointer")
	ErrorInvalidMap          = errors.New("invalid map type, needs to be map[string]any")
)
