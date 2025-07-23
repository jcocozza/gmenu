package renderer

import (
	"github.com/jcocozza/gmenu/renderer/gui"
	"github.com/jcocozza/gmenu/menu"
)

type Renderer interface{
	Init() error
	Cleanup() error
	Render() error
}

func RendererFactory(m menu.Menu) Renderer {
	return &gui.GUIRenderer{M: m}
}
