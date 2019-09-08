package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"path/filepath"
	"strconv"

	"engine/core"
	"engine/graphic"
	"engine/light"

	"engine/loader/gltf"
	"engine/material"
	"engine/math32"
	"engine/util/application"
	"engine/window"

	"github.com/disintegration/imaging"
)

type renderingApp struct {
	application.Application
	x, y, z            float32
	c_imagestream      chan []byte
	c_commands         chan []byte
	Width              int
	Height             int
	jpegQuality        int
	selectionBuffer    map[core.INode][]graphic.GraphicMaterial
	selection_material material.IMaterial
	modelpath          string
}

// setupScene sets up the current scene
func (app *renderingApp) setupScene() {
	app.selection_material = material.NewPhong(math32.NewColor("Red"))
	app.selectionBuffer = make(map[core.INode][]graphic.GraphicMaterial)

	app.Gl().ClearColor(1.0, 1.0, 1.0, 1.0)

	er := app.loadScene(app.modelpath)
	if er != nil {
		log.Fatal(er)
	}

	amb := light.NewAmbient(&math32.Color{0.2, 0.2, 0.2}, 1.0)
	app.Scene().Add(amb)

	plight := light.NewPoint(math32.NewColor("white"), 40)
	plight.SetPosition(100, 20, 70)
	plight.SetLinearDecay(.001)
	plight.SetQuadraticDecay(.001)
	app.Scene().Add(plight)

	app.Camera().GetCamera().SetPosition(1, 1, 3)

	p := math32.Vector3{X: 0, Y: 0, Z: 0}
	app.Camera().GetCamera().LookAt(&p)
	app.Orbit().Enabled = true
	app.Application.Subscribe(application.OnAfterRender, app.onRender)
}

// nameChildren names all gltf nodes by path
func nameChildren(p string, n core.INode) {
	node := n.GetNode()
	node.SetName(p)
	for _, child := range node.Children() {
		idx := node.ChildIndex(child)
		title := p + "/" + strconv.Itoa(idx)
		nameChildren(title, child)
	}
}

type Message struct {
	Action string `json:"action"`
	Value  string `json:"value"`
	Done   bool   `json:"done"`
}

func (a *renderingApp) sendMessageToClient(action string, value string, done bool) {
	m := &Message{Action: action, Value: value, Done: done}
	msg_json, err := json.Marshal(m)
	if err != nil {
		a.Log().Error(err.Error())
		return
	}
	a.Log().Info("sending message: " + string(msg_json))
	a.c_imagestream <- []byte(string(msg_json))
}

// loadScene loads a gltf file
func (a *renderingApp) loadScene(fpath string) error {
	a.sendMessageToClient("loading", fpath, false)
	// Checks file extension
	ext := filepath.Ext(fpath)
	var g *gltf.GLTF
	var err error

	// Parses file
	if ext == ".gltf" {
		g, err = gltf.ParseJSON(fpath)
	} else if ext == ".glb" {
		g, err = gltf.ParseBin(fpath)
	} else {
		return fmt.Errorf("unrecognized file extension:%s", ext)
	}

	if err != nil {
		return err
	}

	defaultSceneIdx := 0
	if g.Scene != nil {
		defaultSceneIdx = *g.Scene
	}

	// Create default scene
	n, err := g.LoadScene(defaultSceneIdx)
	if err != nil {
		return err
	}

	a.Scene().Add(n)
	root := a.Scene().ChildIndex(n)
	nameChildren("/"+strconv.Itoa(root), n)
	a.sendMessageToClient("loading", fpath, true)
	return nil
}

func (a *renderingApp) onRender(evname string, ev interface{}) {
	a.makeScreenShot()
}

// makeScreenShot reads the opengl buffer, encodes it as jpeg and sends it to the channel
func (app *renderingApp) makeScreenShot() {
	w := app.Width
	h := app.Height
	data := app.Gl().ReadPixels(0, 0, w, h, 6408, 5121)
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	img.Pix = data
	img = imaging.FlipV(img)
	buf := new(bytes.Buffer)
	var opt jpeg.Options
	opt.Quality = app.jpegQuality
	jpeg.Encode(buf, img, &opt)
	imageBit := buf.Bytes()
	imgBase64Str := base64.StdEncoding.EncodeToString([]byte(imageBit))
	app.c_imagestream <- []byte(imgBase64Str)
}

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

// selectNode uses a raycaster to get the selected node.
// It sends the selection as json to the image channel
// and changes the node's material
func (app *renderingApp) selectNode(mx float32, my float32) {
	width, height := app.Window().Size()
	x := (-.5 + mx/float32(width)) * 2.0
	y := (.5 - my/float32(height)) * 2.0
	app.Log().Info("click : %f, %f", x, y)
	r := core.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	app.CameraPersp().SetRaycaster(r, x, y)

	i := r.IntersectObject(app.Scene(), true)

	var object *core.Node
	if len(i) != 0 {
		object = i[0].Object.GetNode()

		app.Log().Info("selected : %s", object.Name())
		app.sendMessageToClient("selected", object.Name(), true)
		app.changeNodeMaterial(i[0].Object)
	}
}

// resetSelection resets selected nodes to their original state
func (app *renderingApp) resetSelection() {
	for inode, materials := range app.selectionBuffer {
		gnode, _ := inode.(graphic.IGraphic)
		gfx := gnode.GetGraphic()
		gfx.ClearMaterials()
		for _, material := range materials {
			gfx.AddMaterial(material.IGraphic(), material.IMaterial(), 0, 0)
		}
		delete(app.selectionBuffer, inode)
	}
}

// changeNodeMaterial changes a node's material to selected
func (app *renderingApp) changeNodeMaterial(inode core.INode) {
	gnode, ok := inode.(graphic.IGraphic)
	app.resetSelection()

	if ok {
		if gnode.Renderable() {
			gfx := gnode.GetGraphic()
			var materials []graphic.GraphicMaterial
			for _, material := range gfx.Materials() {
				materials = append(materials, material)
			}
			app.selectionBuffer[inode] = materials
			gfx.ClearMaterials()
			gfx.AddMaterial(gnode, app.selection_material, 0, 0)
		}
	}
}

// commandLoop listens for incoming commands and forwards them to the rendering app
func (app *renderingApp) commandLoop() {
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
