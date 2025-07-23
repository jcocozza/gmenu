package main


type Selector struct {
	curr int
	max int
}

func NewSelector(max int) *Selector {
	return &Selector{
		max: max,
	}
}

func (s *Selector) Curr() int {
	return s.curr
}

func (s *Selector) Left() {
	if s.curr > 0 {
		s.curr--
	}
}
func (s *Selector) Right() {
	if s.curr < s.max - 1 {
		s.curr++
	}
}
