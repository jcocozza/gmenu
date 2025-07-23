package gui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"

	fnt "github.com/jcocozza/gmenu/font"
	"github.com/jcocozza/gmenu/menu"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

type GUIRenderer struct {
	w *glfw.Window

	cursorX float32
	cursorY float32

	width  int
	height int

	M menu.Menu
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

	window, err := glfw.CreateWindow(width, height, "gmenu", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	window.SetPos(0, 0)
	window.MakeContextCurrent()

	window.SetCharCallback(r.charCallback)
	window.SetKeyCallback(r.keyCallback)
	r.w = window
	return nil
}

func (r *GUIRenderer) Cleanup() error {
	glfw.Terminate()
	return nil
}

func (r *GUIRenderer) renderInput(font *glfont.Font) error {
	const inputSize int = 20
	input := r.M.Input()

	var inputDisplay string
	if len(input) < inputSize {
		inputDisplay = input + strings.Repeat(" ", inputSize-len(input))
	} else {
		// this gives us a "scroll" when typing long searches
		start := 0
		end := inputSize
		shift := len(input) - inputSize
		inputDisplay = input[start+shift : end+shift]
	}
	//spacing := strings.Repeat(" ", minSpace)

	font.SetColor(1.0, 1.0, 1.0, 1.0)
	return font.Printf(0, float32(r.height)/2, 1.0, inputDisplay)
}

// TODO: this isn't right at all
func (r *GUIRenderer) renderItems(font *glfont.Font) error {
	items := r.M.Search()
	maxLen := 100
	//minLen := 50

	totalLen := 0

	start := r.M.Current()
	i := start
	for j := start; j < len(items)-1; j++ {
		totalLen += len(items[j]) + 1
		if totalLen > maxLen {
			break
		}
		i++
	}
	displayedList := items[start:i]

	displayStr := /*strings.Repeat(" ", maxLen-totalLen) +*/ strings.Join(displayedList, " ")

	return font.Printf(21, float32(r.height)/2, 1.0, displayStr)
}

func (r *GUIRenderer) Render() error {
	fontSize := int32(float32(r.height) * .6)
	font, err := glfont.LoadFontBytes(fnt.RobotoTTF, fontSize, r.width, r.height)
	if err != nil {
		return err
	}
	for !r.w.ShouldClose() {
		// background color
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		err := r.renderInput(font)
		if err != nil {
			return err
		}
		err = r.renderItems(font)
		if err != nil {
			return err
		}

		r.w.SwapBuffers()
		glfw.PollEvents()
	}
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
		fmt.Println("enter is pressed")
	case glfw.KeyBackspace:
		fmt.Println("backspace")
		r.M.Remove()
	case glfw.KeyUp, glfw.KeyRight:
		fmt.Println("up or right")
	case glfw.KeyDown, glfw.KeyLeft:
		fmt.Println("down or left")
	}
}

func (r *GUIRenderer) charCallback(w *glfw.Window, c rune) {
	r.M.Add(c)
}

//func (r *GUIRenderer) drawCursor(x, y, height float32) {
//	gl.Begin(gl.LINES)
//	gl.Color3f(1.0, 1.0, 1.0) // white
//	gl.Vertex2f(x, y)
//	gl.Vertex2f(x, y+height)
//	gl.End()
//}
//
//func (r *GUIRenderer) drawCursorBlock(x, y, width, height float32) {
//	gl.Begin(gl.QUADS)
//	gl.Color3f(1.0, 1.0, 1.0) // white
//	gl.Vertex2f(x, y)
//	gl.Vertex2f(x+width, y)
//	gl.Vertex2f(x+width, y+height)
//	gl.Vertex2f(x, y+height)
//	gl.End()
//}
