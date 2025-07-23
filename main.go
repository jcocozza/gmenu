package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
	"log"
	"runtime"

	fnt "github.com/jcocozza/gmenu/font"
)

func init() {
	runtime.LockOSThread()
}

var UseStrictCoreProfile = (runtime.GOOS == "darwin")

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)

	if UseStrictCoreProfile {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) // needed on macOS
	}


	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	width := mode.Width
	height := 40

	window, err := glfw.CreateWindow(width, height, "gmenu", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	window.SetPos(0, 0)
	window.MakeContextCurrent()

	font, err := glfont.LoadFontBytes(fnt.RobotoTTF, int32(10), width, height)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		font.SetColor(1.0, 1.0, 1.0, 1.0) // white color
		err := font.Printf(0, float32(height)/2, 1.0, "Hello, GLFW!")
		if err != nil {
			panic(err)
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
