package menu

import "strings"

// a searcher searches
type Searcher interface {
	// return the entire set of searched items if input is empty
	Search(input string) []Item
}

type ListSearcher struct {
	items []Item
}

func (l *ListSearcher) Search(input string) []Item {
	if len(input) == 0 {
		return l.items
	}
	results := []Item{}
	for _, elm := range l.items {
		if strings.Contains(elm.Name(), input) {
			results = append(results, elm)
		}
	}
	return results
}
