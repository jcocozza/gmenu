package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
	"log"
	"runtime"

	fnt "github.com/jcocozza/gmenu/font"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

type GUIRenderer struct {
	w *glfw.Window

	width  int
	height int
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
	height := 40

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

	r.w = window
	return nil
}

func (r *GUIRenderer) Cleanup() error {
	glfw.Terminate()
	return nil
}

func (r *GUIRenderer) Render() {
	font, err := glfont.LoadFontBytes(fnt.RobotoTTF, int32(10), r.width, r.height)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}
	for !r.w.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		font.SetColor(1.0, 1.0, 1.0, 1.0) // white color
		err := font.Printf(0, float32(r.height)/2, 1.0, "Hello, GLFW!")
		if err != nil {
			panic(err)
		}
		r.w.SwapBuffers()
		glfw.PollEvents()
	}
}
