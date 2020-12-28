package renderer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/g3n/engine/window"
)

// Command received from client
type Command struct {
	X     float32
	Y     float32
	Cmd   string
	Val   string
	Moved bool
	Ctrl  bool
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
	t := reflect.TypeOf(app)
	v := reflect.ValueOf(app)
	k := reflect.TypeOf(Command{}).Kind()
	for {
		message := <-app.cCommands

		// retrieve command data from payload
		cmd := Command{}
		err := json.Unmarshal(message, &cmd)
		if err != nil {
			app.Log().Error(err.Error())
		}

		// no command should be directed to orbit control
		if cmd.Cmd == "" {
			cmd.Cmd = "Navigate"
		} else {
			app.Log().Info("received command: %v", cmd)
		}

		// if a func with a matching command name exists,
		// call it with two args: the app itself and the command payload
		m, found := t.MethodByName(cmd.Cmd)
		if found {
			// make sure we got the right func with Command argument
			// otherwise Func.Call will panic
			if m.Type.NumIn() == 2 {
				if m.Type.In(1).Kind() == k {
					args := []reflect.Value{v, reflect.ValueOf(cmd)}
					m.Func.Call(args)
				}
			}
		} else {
			app.Log().Info("Unknown Command: " + cmd.Cmd)
		}
	}
}

// Navigate orbit navigation
func (app *RenderingApp) Navigate(cmd Command) {
	cev := window.CursorEvent{Xpos: cmd.X, Ypos: cmd.Y}
	app.Orbit().OnCursorPos(&cev)
}

// Mousedown triggers a mousedown event
func (app *RenderingApp) Mousedown(cmd Command) {
	mev := window.MouseEvent{Xpos: cmd.X, Ypos: cmd.Y,
		Action: window.Press,
		Button: mapMouseButton(cmd.Val)}
	if cmd.Moved {
		app.imageSettings.isNavigating = true
	}
	app.Orbit().OnMouse(&mev)
}

// Zoom in/out scene
func (app *RenderingApp) Zoom(cmd Command) {
	scrollFactor := float32(10.0)
	mev := window.ScrollEvent{Xoffset: cmd.X, Yoffset: -cmd.Y / scrollFactor}
	app.Orbit().OnScroll(&mev)
}

// Mouseup event
func (app *RenderingApp) Mouseup(cmd Command) {
	mev := window.MouseEvent{Xpos: cmd.X, Ypos: cmd.Y,
		Action: window.Release,
		Button: mapMouseButton(cmd.Val)}

	app.imageSettings.isNavigating = false
	app.Orbit().OnMouse(&mev)

	// mouse left click
	if cmd.Val == "0" && !cmd.Moved {
		app.selectNode(cmd.X, cmd.Y, cmd.Ctrl)
	}
}

// Hide selected element
func (app *RenderingApp) Hide(cmd Command) {
	for inode := range app.selectionBuffer {
		inode.GetNode().SetVisible(false)
	}
	app.resetSelection()
}

// Unhide all hidden elements
func (app *RenderingApp) Unhide(cmd Command) {
	for _, node := range app.nodeBuffer {
		node.SetVisible(true)
	}
}

// Send element userdata to client
func (app *RenderingApp) Userdata(cmd Command) {
	if node, ok := app.nodeBuffer[cmd.Val]; ok {
		app.sendMessageToClient("userdata", fmt.Sprintf("%v", node.UserData()))
	}
}

// Keydown event
func (app *RenderingApp) Keydown(cmd Command) {
	kev := window.KeyEvent{Action: window.Press, Mods: 0, Keycode: mapKey(cmd.Val)}
	app.Orbit().OnKey(&kev)
}

// Keyup event
func (app *RenderingApp) Keyup(cmd Command) {
	kev := window.KeyEvent{Action: window.Release, Mods: 0, Keycode: mapKey(cmd.Val)}
	app.Orbit().OnKey(&kev)
}

// View sets standard views
func (app *RenderingApp) View(cmd Command) {
	app.setCamera(cmd.Val)
}

// Zoomextent entire model
func (app *RenderingApp) Zoomextent(cmd Command) {
	app.zoomToExtent()
}

// Focus on selection
func (app *RenderingApp) Focus(cmd Command) {
	app.focusOnSelection()
}

// Invert image
func (app *RenderingApp) Invert(cmd Command) {
	if app.imageSettings.invert {
		app.imageSettings.invert = false
	} else {
		app.imageSettings.invert = true
	}
}

// Imagesettings applies rendering settings
func (app *RenderingApp) Imagesettings(cmd Command) {
	s := strings.Split(cmd.Val, ":")
	if len(s) == 5 {
		br, err := strconv.Atoi(s[0])
		if err == nil {
			app.imageSettings.brightness = float64(getValueInRange(br, -100, 100))
		}
		ct, err := strconv.Atoi(s[1])
		if err == nil {
			app.imageSettings.contrast = float64(getValueInRange(ct, -100, 100))
		}
		sa, err := strconv.Atoi(s[2])
		if err == nil {
			app.imageSettings.saturation = float64(getValueInRange(sa, -100, 100))
		}
		bl, err := strconv.Atoi(s[3])
		if err == nil {
			app.imageSettings.blur = float64(getValueInRange(bl, 0, 20))
		}
		pix, err := strconv.ParseFloat(s[4], 64)
		if err == nil {
			app.imageSettings.pixelation = getFloatValueInRange(pix, 1.0, 10.0)
		}
	}
}

// Quality settings
func (app *RenderingApp) Quality(cmd Command) {
	quality, err := strconv.Atoi(cmd.Val)
	if err == nil {
		switch quality {
		case 0:
			app.imageSettings.quality = highQ
		case 2:
			app.imageSettings.quality = lowQ
		default:
			app.imageSettings.quality = mediumQ
		}
	}
}

// Enocder settings
func (app *RenderingApp) Encoder(cmd Command) {
	app.imageSettings.encoder = cmd.Val
}

// Fov applies field of view
func (app *RenderingApp) Fov(cmd Command) {
	fov, err := strconv.Atoi(cmd.Val)
	if err == nil {
		app.CameraPersp().SetFov(float32(getValueInRange(fov, 5, 120)))
	}
}

// Debugmode toggles bytegraph
func (app *RenderingApp) Debugmode(cmd Command) {
	if app.Debug {
		app.Debug = false
	} else {
		app.Debug = true
	}
}

// Closes connection
func (app *RenderingApp) Close(cmd Command) {
	app.Log().Info("close")
	app.Window().SetShouldClose(true)
}

// getValueInRange returns a value within bounds
func getValueInRange(value int, lower int, upper int) int {
	if value > upper {
		return upper
	} else if value < lower {
		return lower
	} else {
		return value
	}
}

// getFloatValueInRange returns a value within bounds
func getFloatValueInRange(value float64, lower float64, upper float64) float64 {
	if value > upper {
		return upper
	} else if value < lower {
		return lower
	} else {
		return value
	}
}
