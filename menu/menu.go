package menu

import "strings"

type Menu interface {
	Search() []string
	Add(c rune)
	Remove()
	Input() string
	Current() int
	Left()
	Right()
	//Items() []string
}

func MenuFactory(items []string) Menu {
	return &BasicMenu{
		items: items,
		input: make([]rune, 0),
	}
}

type BasicMenu struct {
	items   []string
	input   []rune
	current int
}

func (s *BasicMenu) Search() []string {
	results := []string{}
	for _, elm := range s.items {
		if strings.Contains(elm, string(s.input)) {
			results = append(results, elm)
		}
	}
	return results
}

func (s *BasicMenu) Add(c rune) {
	s.input = append(s.input, c)
}

func (s *BasicMenu) Remove() {
	if len(s.input) > 0 {
		s.input = s.input[:len(s.input)-1]
	}
}

func (s *BasicMenu) Input() string {
	return string(s.input)
}

func (s *BasicMenu) Current() int {
	return s.current
}

func (s *BasicMenu) Left() {
	if s.current > 0 {
		s.current--
	}
}
func (s *BasicMenu) Right() {
	if s.current < len(s.items)-1 {
		s.current++
	}
}

func (s *BasicMenu) Items() []string {
	return s.items
}
