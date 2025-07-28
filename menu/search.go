package menu

import "strings"

// a searcher searches
type Searcher interface {
	// return the entire set of searched items if input is empty
	Search(input string) []string
}

type ListSearcher struct {
	items []string
}

func (l *ListSearcher) Search(input string) []string {
	if len(input) == 0 {
		return l.items
	}
	results := []string{}
	for _, elm := range l.items {
		if strings.Contains(elm, input) {
			results = append(results, elm)
		}
	}
	return results
}
