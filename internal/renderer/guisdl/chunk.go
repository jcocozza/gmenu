package guisdl

import (
	"fmt"

	"github.com/jcocozza/gmenu/internal/menu"
)

// a chunk is a start and end index in a list
// start and end are exact list indices
type chunk struct {
	start int
	end   int
}

// for debugging
func (c chunk) String() string {
	return fmt.Sprintf("{start: %d end: %d}", c.start, c.end)
}

func (s *SDLRenderer) spaceWidth() int {
	w, _, _ := s.textMetrics("    ")
	return w
}

func (s *SDLRenderer) resultsWidth() float32 {
	return float32(s.width) * (s.searchResultRatio)
}

func (s *SDLRenderer) inputWidth() float32 {
	return float32(s.width) - (float32(s.width) * s.searchResultRatio)
}

func (s *SDLRenderer) chunkify(items []menu.Item, currItemIdx int) ([]chunk, int) {
	spaceWidth := s.spaceWidth()
	maxWidth := s.resultsWidth()

	var chunks []chunk
	var chunkWidth float32
	var start, chunkIdx int
	foundChunk := false

	for i, item := range items {
		display := item.Display()
		if i == currItemIdx {
			display = fmt.Sprintf("[%s]", display)
		}

		itemWidth, _, err := s.textMetrics(display)
		if err != nil {
			panic(err)
		}
		itemWidth += spaceWidth

		if chunkWidth+float32(itemWidth) > maxWidth && chunkWidth > 0 {
			chunks = append(chunks, chunk{start: start, end: i - 1})
			if !foundChunk && currItemIdx >= start && currItemIdx <= i-1 {
				chunkIdx = len(chunks) - 1
				foundChunk = true
			}
			start = i
			chunkWidth = 0
		}
		chunkWidth += float32(itemWidth)
	}

	// Add final chunk
	chunks = append(chunks, chunk{start: start, end: len(items) - 1})
	if !foundChunk {
		chunkIdx = len(chunks) - 1
	}

	return chunks, chunkIdx
}
