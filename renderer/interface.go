package renderer

import (
	"github.com/jcocozza/gmenu/menu"
	"github.com/jcocozza/gmenu/renderer/gui"
)

type Renderer interface{
	// to be called before any attempt at rendering is made
	Init() error
	// defered to cleanup after finished rendering
	Cleanup() error
	// return true if done rendering
	Done() bool
	// render a single frame
	RenderFrame(gm *menu.GMenu) error
	// return the current action
	Action() menu.Action
	// clear the action after it has been processed
	ClearAction()
}

func RendererFactory() Renderer {
	return &gui.GUIRenderer{}
}
