package peak

type Modifier interface {
	Type() string
}

type modifier struct {
	modifierType string
}

const (
	modifierOr    = "or"
	modifierAfter = "after"
)

func (m *modifier) Type() string {
	return m.modifierType
}

func Or() Modifier {
	return &modifier{modifierType: modifierOr}
}

func After() Modifier {
	return &modifier{modifierType: modifierAfter}
}
