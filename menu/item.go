package menu

type Item interface {
	// value to be displayed when rendered
	Display() string
	// actual value of the item - printed when selected
	Value() string
}

type BasicItem string

func (b BasicItem) Display() string  { return string(b) }
func (b BasicItem) Value() string { return string(b) }

type AliasItem struct {
	alias string
	value string
}

func (i *AliasItem) Display() string  { return i.alias }
func (i *AliasItem) Value() string { return i.value }
