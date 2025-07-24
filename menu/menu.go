package menu

import (
	"strings"
)

type Menu interface {
	Search()
	Results() []string
	Add(c rune)
	Remove()
	Input() string
	CurrentIdx() int
	Current() string
	Left()
	Right()

	//Items() []string
}

func MenuFactory(items []string) Menu {
	return &BasicMenu{
		items:  items,
		input:  make([]rune, 0),
		result: items,
	}
}

type BasicMenu struct {
	items   []string
	result  []string
	input   []rune
	current int
}

func (s *BasicMenu) Search() {
	s.current = 0 // always reset current to zero when searching
	results := []string{}
	for _, elm := range s.items {
		if strings.Contains(elm, string(s.input)) {
			results = append(results, elm)
		}
	}
	s.result = results
}

func (s *BasicMenu) Results() []string {
	return s.result
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

func (s *BasicMenu) CurrentIdx() int {
	return s.current
}

func (s *BasicMenu) Current() string {
	return s.result[s.current]
}

func (s *BasicMenu) Left() {
	if s.current > 0 {
		s.current--
	}
}
func (s *BasicMenu) Right() {
	if s.current < len(s.result)-1 {
		s.current++
	}
}

func (s *BasicMenu) Items() []string {
	return s.items
}
