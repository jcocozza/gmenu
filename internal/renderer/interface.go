package renderer

import (
	"github.com/jcocozza/gmenu/internal/menu"
	//"github.com/jcocozza/gmenu/internal/renderer/gui"
	"github.com/jcocozza/gmenu/internal/renderer/guisdl"
)

type Renderer interface {
	// to be called before any attempt at rendering is made
	Init() error
	// defered to cleanup after finished rendering
	Cleanup() error
	// return true if done rendering
	Done() bool
	// populate the action
	PollEvents()
	// render a single frame
	RenderFrame(gm *menu.GMenu) error
	// return the current action
	Action() menu.Action
	// clear the action after it has been processed
	ClearAction()
	// tell renderer that next iteration should close
	MarkClose()
}

func RendererFactory() Renderer {
	//return &gui.GUIRenderer{}
	return &guisdl.SDLRenderer{}
}
