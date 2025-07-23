package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {

	go func() {
		w := new(app.Window)
		w.Option(app.Decorated(false))
		w.Option(app.Title("gmenu"))
		w.Option(app.Size(unit.Dp(1920), unit.Dp(25)))

		if err := draw(w); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func draw(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()
	var ed widget.Editor

	init := false

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if !init {
				init = true
				gtx.Execute(key.FocusCmd{Tag: &ed})
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					med := material.Editor(th, &ed, "search...")
					return med.Layout(gtx)
				}),
			)
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}
