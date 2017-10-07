package goglwindow

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"log"
	"sync"
	"math"
	"fmt"
)

func InitializeWindow2() *glfw.Window {
	return InitializeWindow(1600, 1200, "Test program")
}

func InitializeWindow(width, height int, name string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	fmt.Printf("GLFW initialized!\n")

	major, minor, revision := glfw.GetVersion()
	fmt.Printf("HOLA %v, %v, %v\n", major, minor, revision)
	if major < 3 {
		panic("Unsupported OpenGL version")
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

type Event interface {

}

type CurserMoved struct {
	X, Y float64
	PX, PY float64
	Delta bool
}

type KeyboardEvent struct {
	Action glfw.Action
	Key glfw.Key
}

type MouseEvent struct {
	Button glfw.MouseButton
	Action glfw.Action
}

type MouseScroll struct {
	X, Y float64
}

type Properties struct {
	storeCurser bool
	x, y float64
}

func NewProperties(storeCurser bool) Properties {
	return Properties{storeCurser, 0,0}
}

type Receiver interface {
	Receive() chan Event
}

type Window struct {
	sync.RWMutex
	Window *glfw.Window
	properties Properties
	reciever chan Event
}

func NewWindow2(keyboard, mouse Receiver) *Window {
	return NewWindow(1600, 1200, "Test program", keyboard, mouse)
}

func NewWindow(width, height int, name string, keyboard, mouse Receiver) *Window {
	return NewWindow2WithProperties(width, height, name, Properties{true, 0,0}, keyboard, mouse)
}

func NewWindow2WithProperties(width, height int, name string, props Properties, keyboard, mouse Receiver) *Window {
	window := &Window{
		Window:InitializeWindow(width, height, name),
		properties:props,
	}

	// Mouse buttons
	window.Window.SetMouseButtonCallback(func(gw *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		mouse.Receive() <- MouseEvent{button, action}
	})

	// Keyboard keys
	window.Window.SetKeyCallback(func(gw *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Repeat {
			return
		}

		keyboard.Receive() <- KeyboardEvent{action, key}
	})

	// Mouse scroll
	window.Window.SetScrollCallback(func(gw *glfw.Window, x, y float64) {
		mouse.Receive() <- MouseScroll{x, y}
	})

	// Mouse position
	window.Window.SetCursorPosCallback(func(gw *glfw.Window, xpos float64, ypos float64) {
		window.Lock()
		defer window.Unlock()

		cm := CurserMoved{}
		cm.PX = xpos
		cm.PY = ypos

		lastX := window.properties.x
		lastY := window.properties.y
		window.properties.x = xpos
		window.properties.y = ypos

		if lastX == math.Inf(-1) && lastY == math.Inf(-1) {
			return
		}

		cm.X = xpos - lastX
		cm.Y = ypos - lastY

		mouse.Receive() <- cm
	})

	return window
}
