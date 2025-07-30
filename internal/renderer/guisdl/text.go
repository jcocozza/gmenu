package guisdl

import "github.com/veandco/go-sdl2/sdl"

// width, height, err
func (s *SDLRenderer) textMetrics(text string) (int, int, error) {
	return s.font.SizeUTF8(text)
}

var (
	White      = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	Hightlight = sdl.Color{R: 255, G: 1.0, B: 1.0, A: 255}
	//Black =
)

func (s *SDLRenderer) drawText(text string, color sdl.Color, x int, y int) error {
	textSurface, err := s.font.RenderUTF8Blended(text, color)
	if err != nil {
		return err
	}
	defer textSurface.Free()

	// Create texture from surface
	textTexture, err := s.renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		return err
	}
	defer textTexture.Destroy()
	textRect := sdl.Rect{X: int32(x), Y: int32(y), W: textSurface.W, H: textSurface.H}
	return s.renderer.Copy(textTexture, nil, &textRect)
}
