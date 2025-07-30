package menu

// must call ExecSearch to populuate the results
type GMenu struct {
	input          []rune
	currentResults []Item
	selected       int

	s Searcher
}

func NewGmenu(items []Item) *GMenu {
	return &GMenu{
		s: SearcherFactory(items),
	}
}

func (g *GMenu) AddChar(c rune) {
	g.input = append(g.input, c)
}

func (g *GMenu) RemoveChar() {
	if len(g.input) == 0 { return }
	g.input = g.input[:len(g.input)-1]
}

func (g *GMenu) Input() string {
	return string(g.input)
}

func (g *GMenu) ExecSearch() {
	g.selected = 0
	g.currentResults = g.s.Search(string(g.input))
}

func (g *GMenu) Results() []Item {
	return g.currentResults
}

func (g *GMenu) Left() {
	if g.selected > 0 {
		g.selected--
	}
}

func (g *GMenu) Right() {
	if g.selected < len(g.currentResults)-1 {
		g.selected++
	}
}

func (g *GMenu) Selected() int {
	return g.selected
}

// return nil if there are no results
func (g *GMenu) SelectedItem() Item {
	if len(g.currentResults) == 0 {
		return nil
	}
	return g.currentResults[g.selected]
}
