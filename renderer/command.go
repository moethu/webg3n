package renderer

import (
	"encoding/json"
	"engine/window"
	"fmt"
	"strconv"
)

type Command struct {
	X     float32
	Y     float32
	Cmd   string
	Val   string
	Moved bool
}

// mapMouseButton maps js mouse buttons to window mouse buttons
func mapMouseButton(value string) window.MouseButton {
	switch value {
	case "0":
		return window.MouseButtonLeft
	case "1":
		return window.MouseButtonMiddle
	case "2":
		return window.MouseButtonRight
	default:
		return window.MouseButtonLeft
	}
}

// mapKey maps js keys to window keys
func mapKey(value string) window.Key {
	switch value {
	case "38":
		return window.KeyUp
	case "37":
		return window.KeyLeft
	case "39":
		return window.KeyRight
	case "40":
		return window.KeyDown
	default:
		return window.KeyEnter
	}
}

// commandLoop listens for incoming commands and forwards them to the rendering app
func (app *RenderingApp) commandLoop() {
	for {
		message := <-app.c_commands

		cmd := Command{}
		err := json.Unmarshal(message, &cmd)
		if err != nil {
			app.Log().Error(err.Error())
		}

		switch cmd.Cmd {
		case "":
			cev := window.CursorEvent{Xpos: cmd.X, Ypos: cmd.Y}
			app.Orbit().OnCursorPos(&cev)
		case "mousedown":
			mev := window.MouseEvent{Xpos: cmd.X, Ypos: cmd.Y,
				Action: window.Press,
				Button: mapMouseButton(cmd.Val)}
			app.Orbit().OnMouse(&mev)
		case "zoom":
			mev := window.ScrollEvent{Xoffset: cmd.X, Yoffset: -cmd.Y}
			app.Orbit().OnScroll(&mev)
		case "mouseup":
			mev := window.MouseEvent{Xpos: cmd.X, Ypos: cmd.Y,
				Action: window.Release,
				Button: mapMouseButton(cmd.Val)}
			app.Orbit().OnMouse(&mev)

			// mouse left click
			if cmd.Val == "0" && !cmd.Moved {
				app.selectNode(cmd.X, cmd.Y)
			}
		case "hide":
			if node, ok := app.nodeBuffer[cmd.Val]; ok {
				node.SetVisible(false)
			}
		case "unhide":
			for _, node := range app.nodeBuffer {
				node.SetVisible(true)
			}
		case "userdata":
			if node, ok := app.nodeBuffer[cmd.Val]; ok {
				app.sendMessageToClient("userdata", fmt.Sprintf("%v", node.UserData()))
			}
		case "keydown":
			kev := window.KeyEvent{Action: window.Press, Mods: 0, Keycode: mapKey(cmd.Val)}
			app.Orbit().OnKey(&kev)
		case "keyup":
			kev := window.KeyEvent{Action: window.Release, Mods: 0, Keycode: mapKey(cmd.Val)}
			app.Orbit().OnKey(&kev)
		case "fov":
			fov, err := strconv.Atoi(cmd.Val)
			if err == nil {
				app.CameraPersp().SetFov(float32(fov))
			}
		case "close":
			app.Log().Info("close")
			app.Window().SetShouldClose(true)
		default:
			app.Log().Info("Unknown Command: " + cmd.Cmd)
		}
	}
}
