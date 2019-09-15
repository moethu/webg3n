package renderer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/g3n/engine/window"
)

// Incoming Command from js
type Command struct {
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Cmd   string  `json:"cmd"`
	Val   string  `json:"val"`
	Moved bool    `json:"moved"`
	Ctrl  bool    `json:"ctrl"`
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

// maps view string to enum
func mapView(value string) Standardview {
	switch value {
	case "front":
		return Front
	case "rear":
		return Rear
	case "top":
		return Top
	case "bottom":
		return Bottom
	case "left":
		return Left
	case "right":
		return Right
	default:
		return Top
	}
}

// commandLoop listens for incoming commands and forwards them to the rendering app
func (app *RenderingApp) commandLoop() {
	for {
		message := <-app.c_commands

		cmd := Command{}
		err := json.Unmarshal(message, &cmd)
		if err != nil {
			app.log.Error(err.Error())
		}

		if cmd.Cmd != "" {
			app.log.Info("received command: %v", cmd)
		}

		switch cmd.Cmd {
		case "":
			cev := mouseEvent{X: cmd.X, Y: cmd.Y}
			app.orbit.OnCursorPos(&cev)
		case "mousedown":
			mev := mouseEvent{X: cmd.X, Y: cmd.Y,
				Button: mapMouseButton(cmd.Val), MouseDown: true}
			app.orbit.OnMouse(&mev)
		case "zoom":
			mev := mouseEvent{X: cmd.X, Y: cmd.Y}
			app.orbit.OnScroll(&mev)
		case "mouseup":
			mev := mouseEvent{X: cmd.X, Y: cmd.Y,
				Button: mapMouseButton(cmd.Val), MouseDown: false}
			app.orbit.OnMouse(&mev)

			// mouse left click
			if cmd.Val == "0" && !cmd.Moved {
				app.selectNode(cmd.X, cmd.Y, cmd.Ctrl)
			}
		case "hide":
			for inode, _ := range app.selectionBuffer {
				inode.GetNode().SetVisible(false)
			}
			app.resetSelection()
		case "unhide":
			for _, node := range app.nodeBuffer {
				node.SetVisible(true)
			}
		case "userdata":
			if node, ok := app.nodeBuffer[cmd.Val]; ok {
				app.respondToClient("userdata", fmt.Sprintf("%v", node.UserData()))
			}
		case "keydown":
			kev := keyEvent{Key: mapKey(cmd.Val), IsPressed: true}
			app.orbit.OnKey(&kev)
		case "view":
			app.SetStandardView(mapView(cmd.Val))
		case "zoomextent":
			app.ZoomExtent()
		case "focus":
			app.FocusOnSelection()
		case "invert":
			if app.imageSettings.invert {
				app.imageSettings.invert = false
			} else {
				app.imageSettings.invert = true
			}
		case "imagesettings":
			s := strings.Split(cmd.Val, ":")
			if len(s) == 4 {
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
			}
		case "quality":
			quality, err := strconv.Atoi(cmd.Val)
			if err == nil {
				app.imageSettings.jpegQuality = getValueInRange(quality, 5, 100)
			}
		case "fov":
			fov, err := strconv.Atoi(cmd.Val)
			if err == nil {
				app.camPersp.SetFov(float32(getValueInRange(fov, 5, 120)))
			}
		case "close":
			app.log.Info("close")
			//app.Window().SetShouldClose(true)
		default:
			app.log.Info("Unknown Command: " + cmd.Cmd)
		}
	}
}

// getValueInRange returns an integer value within a certain range
func getValueInRange(value int, lower int, upper int) int {
	if value > upper {
		return upper
	} else if value < lower {
		return lower
	} else {
		return value
	}
}
