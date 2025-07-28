package menu

import (
	"strings"
)

type Item interface {
	// value to be displayed when rendered
	Display() string
	// actual value of the item - printed when selected
	Value() string
}

type BasicItem string

func (b BasicItem) Display() string { return string(b) }
func (b BasicItem) Value() string   { return string(b) }

type AliasItem struct {
	display string
	value   string
}

func (i AliasItem) Display() string { return i.display }
func (i AliasItem) Value() string   { return i.value }

func ParseAliasItem(ln string) AliasItem {
	lns := strings.Split(ln, " ")

	var display string
	var value string

	switch len(lns) {
	case 0:
		display = ""
		value = ""
	case 1:
		display = lns[0]
		value = lns[0]
	default:
		display = lns[0]
		value = strings.Join(lns[1:], " ")
	}

	return AliasItem{
		display: display,
		value:   value,
	}
}
