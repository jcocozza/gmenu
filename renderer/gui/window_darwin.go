//go:build darwin

package gui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import "window_darwin.h"
*/
import "C"
import "github.com/go-gl/glfw/v3.3/glfw"

func setWindowBehavior(window *glfw.Window) {
	windowPtr := window.GetCocoaWindow()
	C.setWindowBehavior(windowPtr)
}
