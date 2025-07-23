package main

import "github.com/jcocozza/gmenu/renderer"

func main() {
	r := renderer.RendererFactory()
	if err := r.Init(); err != nil {
		panic(err)
	}
	defer r.Cleanup()
	r.Render()
}
