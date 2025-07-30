package guisdl

import (
	"fmt"
	"unicode/utf8"

	"github.com/jcocozza/gmenu/internal/menu"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SDLRenderer struct {
	w        *sdl.Window
	font     *ttf.Font
	renderer *sdl.Renderer

	height            int
	width             int
	searchResultRatio float32

	currAction menu.Action

	done bool
}

func (s *SDLRenderer) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	mode, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		return err
	}

	s.width = int(mode.W)
	s.height = 200

	window, err := sdl.CreateWindow("gmenu", 0, 0, int32(s.width), int32(s.height), sdl.WINDOW_BORDERLESS|sdl.WINDOW_ALWAYS_ON_TOP|sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	s.w = window

	window.SetAlwaysOnTop(true)
	window.SetResizable(false)

	if err = ttf.Init(); err != nil {
		return err
	}

	font, err := ttf.OpenFont("/Users/josephcocozza/Repositories/gmenu/font/Roboto-Regular.ttf", 12)
	if err != nil {
		return err
	}
	s.font = font

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	s.renderer = renderer

	sdl.StartTextInput()
	return nil
}
func (s *SDLRenderer) Cleanup() error {
	sdl.Quit()
	s.font.Close()
	ttf.Quit()
	err := s.renderer.Destroy()
	if err != nil {
		return err
	}
	sdl.StopTextInput()
	return s.w.Destroy()
}
func (s *SDLRenderer) Done() bool { return s.done }
func (s *SDLRenderer) PollEvents() {
	switch e := sdl.PollEvent().(type) {
	case *sdl.KeyboardEvent:
		if e.Type != sdl.KEYDOWN {
			return
		}
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			s.done = true
		case sdl.K_RETURN, sdl.K_RETURN2:
			s.currAction = menu.ActionSelect{}
		case sdl.K_DELETE, sdl.K_BACKSPACE:
			s.currAction = menu.ActionRemoveChar{}
		case sdl.K_DOWN, sdl.K_LEFT:
			s.currAction = menu.ActionLeft{}
		case sdl.K_UP, sdl.K_RIGHT:
			s.currAction = menu.ActionRight{}
		}
	case *sdl.TextInputEvent:
		r, _ := utf8.DecodeRune(e.Text[:])
		s.currAction = menu.ActionAddChar{C: r}
	}
}
func (s *SDLRenderer) RenderFrame(gm *menu.GMenu) error {
	if err := s.renderer.SetDrawColor(0, 0, 0, 255); err != nil {
		return err
	} // Black background
	if err := s.renderer.Clear(); err != nil {
		return err
	}
	height := float32(s.height) / 2

	//
	if len(gm.Input()) != 0 {
		if err := s.drawText(gm.Input(), White, 0, int(height)); err != nil {
			return err
		}
	}

	// render the result
	spaceWidth := float32(s.spaceWidth())
	items := gm.Results()
	selected := gm.Selected()

	chunks, chunkIdx := s.chunkify(items, selected)
	offset := float32(s.width) - s.resultsWidth()
	var displayWidth float32
	c := chunks[chunkIdx]
	for i, elm := range items[c.start : c.end+1] {
		isCurrent := c.start+i == selected
		displayItem := elm.Display()

		color := White
		if isCurrent {
			color = Hightlight
			displayItem = fmt.Sprintf("[%s]", displayItem)
		}
		if err := s.drawText(displayItem, color, int(offset+displayWidth), int(height)); err != nil {
			return err
		}
		itemWidth, _, err := s.textMetrics(displayItem)
		if err != nil { return err }
		displayWidth += float32(itemWidth) + spaceWidth
	}

	s.renderer.Present()
	return nil
}
func (s *SDLRenderer) Action() menu.Action { return s.currAction }
func (s *SDLRenderer) ClearAction()        { s.currAction = nil }
func (s *SDLRenderer) MarkClose()          { s.done = true }
