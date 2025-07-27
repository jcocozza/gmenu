package gui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import "window_mac.h"
*/
import "C"
import "github.com/go-gl/glfw/v3.3/glfw"

const (
	NSWindowCollectionBehaviourCanJoinAllSpaces = 1 << 0
	NSWindowCollectionBehaviorFullScreenAuxiliary = 1 << 9
)

func setWindowCollectionBehavior(window *glfw.Window) {
	windowPtr := window.GetCocoaWindow()
	C.setWindowCollectionBehavior(windowPtr)

}
