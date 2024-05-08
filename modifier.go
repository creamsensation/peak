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

func Or(use ...bool) Modifier {
	if len(use) > 0 && !use[0] {
		return nil
	}
	return &modifier{modifierType: modifierOr}
}

func After(use ...bool) Modifier {
	if len(use) > 0 && !use[0] {
		return nil
	}
	return &modifier{modifierType: modifierAfter}
}
