package gui

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

func (r *GUIRenderer) resultsWidth() float32 {
	return float32(r.width) * r.searchResultRatio
}

func (r *GUIRenderer) inputWidth() float32 {
	return float32(r.width) - (float32(r.width) * r.searchResultRatio)
}

func (r *GUIRenderer) spaceWidth() float32 {
	return r.f.Width(r.scale, "    ")
}


func (r *GUIRenderer) chunkify(items []menu.Item, currItemIdx int) ([]chunk, int) {
	spaceWidth := r.spaceWidth()
	maxWidth := r.resultsWidth()
	var chunkWidth float32

	var chunks []chunk
	var start int
	chunkIdx := -1
	for i, item := range items {
		displayItem := item.Display()
		isCurrent := i == currItemIdx
		if isCurrent {
			displayItem = fmt.Sprintf("[%s]", item.Display())
		}

		itemWidth := r.f.Width(r.scale, "%s", displayItem) + spaceWidth

		if chunkWidth+itemWidth >= maxWidth {
			chunks = append(chunks, chunk{start: start, end: i-1})

			if currItemIdx >= start && currItemIdx <= i-1 {
				chunkIdx = len(chunks) - 1
			}
			chunkWidth = 0
			start = i
		}
		chunkWidth += itemWidth
	}

	// add the last chunk
	chunks = append(chunks, chunk{start: start, end: len(items) - 1})
	if chunkIdx == -1 {
		chunkIdx = len(chunks) - 1
	}
	return chunks, chunkIdx
}
