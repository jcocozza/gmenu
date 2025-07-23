package main

import (
	"github.com/jcocozza/gmenu/menu"
	"github.com/jcocozza/gmenu/renderer"
)

func main() {
	m := menu.MenuFactory(words)
	r := renderer.RendererFactory(m)
	if err := r.Init(); err != nil {
		panic(err)
	}
	defer r.Cleanup()
	err := r.Render()
	if err != nil {
		panic(err)
	}
}
