package gui

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jcocozza/gmenu/menu"
	"github.com/nullboundary/glfont"

	fnt "github.com/jcocozza/gmenu/font"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

// MUST call Init()
type GUIRenderer struct {
	w *glfw.Window

	scale             float32
	width             int
	height            int
	searchResultRatio float32

	f *glfont.Font

	currAction menu.Action
}

func (r *GUIRenderer) Init() error {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)

	if useStrictCoreProfile {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) // needed on macOS
	}

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	width := mode.Width
	height := 20

	r.width = width
	r.height = height
	r.scale = 1.0
	r.searchResultRatio = 0.6

	window, err := glfw.CreateWindow(width, height, "gmenu", nil, nil)
	if err != nil {
		return err
	}

	window.SetPos(0, 0)
	window.MakeContextCurrent()
	window.Show()
	window.Focus()
	window.RequestAttention()

	window.SetCharCallback(r.charCallback)
	window.SetKeyCallback(r.keyCallback)
	r.w = window

	setWindowBehavior(window)

	if err := gl.Init(); err != nil {
		return err
	}

	fontSize := int32(float32(r.height) * .6)
	font, err := glfont.LoadFontBytes(fnt.RobotoTTF, fontSize, r.width, r.height)
	if err != nil {
		return err
	}
	r.f = font

	return nil
}

func (r *GUIRenderer) Cleanup() error {
	glfw.Terminate()
	return nil
}

func (r *GUIRenderer) Done() bool {
	return r.w.ShouldClose()
}

func (r *GUIRenderer) RenderFrame(gm *menu.GMenu) error {
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// this renders the input
	input := gm.Input()
	var inputDisplay string
	if len(input) < int(r.inputWidth()) {
		inputDisplay = input
	} else {
		start := 0
		end := int(r.inputWidth())
		shift := len(input) - int(r.inputWidth())
		inputDisplay = input[start+shift : end+shift]
	}
	r.f.SetColor(1.0, 1.0, 1.0, 1.0)
	if err := r.f.Printf(0, float32(r.height)/2, r.scale, "%s", inputDisplay); err != nil {
		return err
	}

	// this renders the results
	items := gm.Results()
	selected := gm.Selected()

	chunks, chunkIdx := r.chunkify(items, selected)
	offset := float32(r.width) - r.resultsWidth()
	height := float32(r.height) / 2

	var displayWidth float32
	c := chunks[chunkIdx]
	for i, elm := range items[c.start : c.end+1] {
		isCurrent := c.start+i == selected
		displayItem := elm.Display()
		if isCurrent {
			r.f.SetColor(255, 1.0, 1.0, 255)
			displayItem = fmt.Sprintf("[%s]", elm)
		}
		if err := r.f.Printf(offset+displayWidth, height, r.scale, "%s", displayItem); err != nil {
			return err
		}
		if isCurrent {
			r.f.SetColor(1.0, 1.0, 1.0, 1.0)
		}
		itemWidth := r.f.Width(r.scale, "%s", displayItem)
		displayWidth += itemWidth + r.spaceWidth()
	}

	r.w.SwapBuffers()
	return nil
}

func (r *GUIRenderer) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press && action != glfw.Repeat {
		return
	}
	switch key {
	case glfw.KeyEscape:
		w.SetShouldClose(true)
	case glfw.KeyEnter:
		r.currAction = menu.ActionSelect{}
		w.SetShouldClose(true)
	case glfw.KeyBackspace:
		r.currAction = menu.ActionRemoveChar{}
	case glfw.KeyUp, glfw.KeyRight:
		r.currAction = menu.ActionRight{}
	case glfw.KeyDown, glfw.KeyLeft:
		r.currAction = menu.ActionLeft{}
	}
}

func (r *GUIRenderer) charCallback(w *glfw.Window, c rune) {
	r.currAction = menu.ActionAddChar{C: c}
}

func (r *GUIRenderer) Action() menu.Action { return r.currAction }
func (r *GUIRenderer) ClearAction()        { r.currAction = nil }
