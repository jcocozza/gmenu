package main

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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

func drawRectWithText(gtx layout.Context, lbl material.LabelStyle, bgColor color.NRGBA) layout.Dimensions {
	dims := lbl.Layout(gtx)
	paint.FillShape(gtx.Ops, bgColor, clip.Rect{Max: dims.Size}.Op())
	lbl.Layout(gtx)
	return dims
}

func draw(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()
	var ed widget.Editor
	ed.SingleLine = true
	var list widget.List
	list.Axis = layout.Horizontal

	init := false
	sl := NewSearchList(words)
	selector := NewSelector(len(words))

	for {
		filteredWords := sl.Search(ed.Text())

		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if !init {
				init = true
				gtx.Execute(key.FocusCmd{Tag: &ed})
			}

			pressedKeys := make(map[key.Name]bool)
			for {
				event, ok := gtx.Event(key.Filter{Name: key.NameEscape}, key.Filter{Name: key.NameEnter}, key.Filter{Name: key.NameLeftArrow}, key.Filter{Name: key.NameRightArrow})
				if !ok {
					break
				}
				switch event := event.(type) {
				case key.Event:
					if event.State == key.Press {
						if pressedKeys[event.Name] {
							continue
						}
						pressedKeys[event.Name] = true
						switch event.Name {
						case key.NameEscape:
							return nil
						case key.NameEnter:
							return nil
						case key.NameRightArrow, key.NameUpArrow:
							selector.Right()
						case key.NameLeftArrow, key.NameDownArrow:
							selector.Left()
						}
					} else if event.State == key.Release {
						pressedKeys[event.Name] = false
					}
				}
			}

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					med := material.Editor(th, &ed, "search...")
					return med.Layout(gtx)
				}),
				// Spacer that expands to take all remaining space
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Max} // Just occupy space, no drawing needed
				}),

				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					list.ScrollTo(selector.curr)
					return list.Layout(gtx, len(filteredWords), func(gtx layout.Context, i int) layout.Dimensions {
						med := material.Body1(th, filteredWords[i])
						if i == selector.curr {
							med.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
							background := color.NRGBA{R: 50, G: 100, B: 200, A: 255}
							return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawRectWithText(gtx, med, background)
							})
						}
						return layout.UniformInset(unit.Dp(4)).Layout(gtx, med.Layout)
					})
				}),
			)
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}
