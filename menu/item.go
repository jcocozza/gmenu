package menu

type Item interface {
	Name() string
	Value() string
	Equals(other string) bool
}

type BasicItem string

func (b BasicItem) Name() string  { return string(b) }
func (b BasicItem) Value() string { return string(b) }

type AliasItem struct {
	alias string
	value string
}

func (i *AliasItem) Name() string  { return i.alias }
func (i *AliasItem) Value() string { return i.value }
