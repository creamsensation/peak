package peak

type fragmentBuilder struct {
	builders []QueryBuilder
}

func Fragment(builders ...QueryBuilder) QueryBuilder {
	return &fragmentBuilder{
		builders: builders,
	}
}

func (b *fragmentBuilder) Build() BuildResult {
	return BuildResult{}
}
