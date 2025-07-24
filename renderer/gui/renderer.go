package gui

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"

	fnt "github.com/jcocozza/gmenu/font"
	"github.com/jcocozza/gmenu/menu"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

const scale float32 = 1

type GUIRenderer struct {
	w *glfw.Window

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

	if err := gl.Init(); err != nil {
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
	font.SetColor(1.0, 1.0, 1.0, 1.0)
	return font.Printf(0, float32(r.height)/2, scale, "%s", inputDisplay)
}

type chunk struct {
	start int
	end   int
}

func (c chunk) String() string {
	return fmt.Sprintf("{start: %d end: %d}", c.start, c.end)
}

func (r *GUIRenderer) rItems(font *glfont.Font) error {
	maxWidth := (float32(r.width) / 2 ) - 4 * font.Width(scale, "    ")// give the "items" half the window
	items := r.M.Results()
	curr := r.M.CurrentIdx()
	displayChunks := []chunk{}

	start := 0
	end := 0
	currentChunkIdx := 0
	spaceWidth := font.Width(scale, "    ")

	var chunkWidth float32 = 0.0
	for i, item := range items {
		displayItem := item
		isCurrent := i == curr
		if isCurrent {
			displayItem = fmt.Sprintf("[%s]", item) // cute highlight
		}
		itemWidth := font.Width(scale, "%s", displayItem) + spaceWidth
		if chunkWidth+itemWidth >= maxWidth { // we have gone too far
			displayChunks = append(displayChunks, chunk{start: start, end: i})
			if curr >= start && curr <= i {
				currentChunkIdx = len(displayChunks) - 1
			}
			start = i
			end = i
			chunkWidth = itemWidth
		} else {
			chunkWidth += itemWidth
		}
	}
	// add the last chunk
	displayChunks = append(displayChunks, chunk{start: end, end: len(items)})

	// actual render
	offset := float32(r.width) / 2 // start halfway across the screen
	height := float32(r.height) / 2
	var displayWidth float32 = 0.0
	c := displayChunks[currentChunkIdx]
	for i, elm := range items[c.start:c.end] {
		isCurrent := c.start+i == curr
		displayItem := elm
		if isCurrent {
			font.SetColor(255, 1.0, 1.0, 255)
		}
		if isCurrent {
			displayItem = fmt.Sprintf("[%s]", elm) // cute highlight
		}
		if err := font.Printf(offset+displayWidth, height, scale, "%s", displayItem); err != nil {
			return err
		}
		if isCurrent {
			font.SetColor(1.0, 1.0, 1.0, 1.0)
		}
		itemWidth := font.Width(scale, "%s", displayItem)
		displayWidth += itemWidth + spaceWidth
	}
	return nil
}

func (r *GUIRenderer) Render() error {
	fontSize := int32(float32(r.height) * .6)
	font, err := glfont.LoadFontBytes(fnt.RobotoTTF, fontSize, r.width, r.height)
	if err != nil {
		return err
	}
	for !r.w.ShouldClose() {
		glfw.PollEvents()
		// background color
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		err := r.renderInput(font)
		if err != nil {
			return err
		}
		err = r.rItems(font)
		if err != nil {
			return err
		}

		r.w.SwapBuffers()
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
		fmt.Fprintln(os.Stdout, r.M.Current())
		w.SetShouldClose(true)
	case glfw.KeyBackspace:
		r.M.Remove()
		r.M.Search()
	case glfw.KeyUp, glfw.KeyRight:
		r.M.Right()
	case glfw.KeyDown, glfw.KeyLeft:
		r.M.Left()
	}
}

func (r *GUIRenderer) charCallback(w *glfw.Window, c rune) {
	r.M.Add(c)
	r.M.Search()
}
