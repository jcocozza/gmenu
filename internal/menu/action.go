package menu

type Action interface {
	Apply(m *GMenu)
}

type ActionAddChar struct {
	C rune
}

func (a ActionAddChar) Apply(m *GMenu) { m.AddChar(a.C); m.ExecSearch() }

type ActionRemoveChar struct{}

func (a ActionRemoveChar) Apply(m *GMenu) { m.RemoveChar(); m.ExecSearch() }

type ActionLeft struct{}

func (a ActionLeft) Apply(m *GMenu) { m.Left() }

type ActionRight struct{}

func (a ActionRight) Apply(m *GMenu) { m.Right() }

type ActionSelect struct{}

func (a ActionSelect) Apply(m *GMenu) {}
