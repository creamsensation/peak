package peak

func If(condition bool, qb QueryBuilder) QueryBuilder {
	if !condition {
		return nil
	}
	return qb
}
