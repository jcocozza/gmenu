package main

import "strings"

type SearchList struct {
	elements []string
}

func NewSearchList(elms []string) *SearchList {
	return &SearchList{elements: elms}
}

// sample search algorithm
//
// very slow, very bad
func (s *SearchList) Search(input string) []string {
	if len(input) == 0 {
		return s.elements
	}
	l := []string{}
	for _, elm := range s.elements {
		if strings.Contains(elm, input) {
			l = append(l, elm)
		}
	}
	return l
}
