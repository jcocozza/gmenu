package menu

import "strings"

// a searcher searches
type Searcher interface {
	// return the entire set of searched items if input is empty
	Search(input string, ignoreCase bool) []Item
}

func SearcherFactory(items []Item) Searcher {
	return &ListSearcher{
		items: items,
	}
}

type ListSearcher struct {
	items []Item
}

func (l *ListSearcher) Search(input string, ignoreCase bool) []Item {
	if len(input) == 0 {
		return l.items
	}
	results := []Item{}

	if ignoreCase {
		for _, elm := range l.items {
			if strings.Contains(strings.ToLower(elm.Display()), strings.ToLower(input)) {
				results = append(results, elm)
			}
		}
	} else {
		for _, elm := range l.items {
			if strings.Contains(elm.Display(), input) {
				results = append(results, elm)
			}
		}
	}
	return results
}
