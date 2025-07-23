package menu

import "strings"

type Menu interface {
	Search() []string
	Add(c rune)
	Remove()
	Input() string
}

func MenuFactory() Menu {
	return &BasicMenu{
		items: make([]string, 0),
		input: make([]rune, 0),
	}
}

type BasicMenu struct {
	items []string
	input []rune
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
