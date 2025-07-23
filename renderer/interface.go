package renderer

import "github.com/jcocozza/gmenu/renderer/gui"

type Renderer interface{
	Init() error
	Cleanup() error
	Render()
}

func RendererFactory() Renderer {
	return &gui.GUIRenderer{}
}
