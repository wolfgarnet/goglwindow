package goglwindow

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"fmt"
)

type State uint8

const (
	Released State = iota
	Pressed
	WasPressed
	WasReleased
)

type Keyboard struct {
	Keys map[glfw.Key]State
	events         chan Event

	currentKeys, lastKeys []glfw.Key
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		events:make(chan Event, 256),
		Keys:make(map[glfw.Key]State),
	}
}

func (k *Keyboard) WasPressed(key glfw.Key) bool {
	return k.Keys[key] == WasPressed
}

func (k *Keyboard) IsPressed(key glfw.Key) bool {
	return k.Keys[key] == Pressed || k.Keys[key] == WasPressed
}

func (k *Keyboard) IsPressed2(key1, key2 glfw.Key) bool {
	return (k.Keys[key1] == Pressed || k.Keys[key1] == WasPressed) && (k.Keys[key2] == Pressed || k.Keys[key2] == WasPressed)
}

func (k *Keyboard) WasReleased(key glfw.Key) bool {
	return k.Keys[key] == WasReleased
}

func (k *Keyboard) Receive() chan Event {
	return k.events
}

func (k *Keyboard) Consume() {
	k.currentKeys = nil
	done := false
	for !done {
		select {
		case event, ok := <-k.events:
			if ok {
				ke := event.(KeyboardEvent)
				if ke.Action == glfw.Press {
					k.Keys[ke.Key] = WasPressed
				} else {
					k.Keys[ke.Key] = WasReleased
				}

				k.currentKeys = append(k.currentKeys, ke.Key)
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			done = true
		}
	}

	// Check last keys
	for _, key := range k.lastKeys {
		switch k.Keys[key] {
		case WasReleased:
			k.Keys[key] = Released

		case WasPressed:
			k.Keys[key] = Pressed
		}
	}

	k.lastKeys = k.currentKeys
}
