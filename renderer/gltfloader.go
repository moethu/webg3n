package renderer

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/loader/gltf"
	"path/filepath"
	"strconv"
)

// nameChildren names all gltf nodes by path
func (app *RenderingApp) nameChildren(p string, n core.INode) {
	node := n.GetNode()
	node.SetName(p)
	app.nodeBuffer[p] = node
	for _, child := range node.Children() {
		idx := node.ChildIndex(child)
		title := p + "/" + strconv.Itoa(idx)
		app.nameChildren(title, child)
	}
}

// loadScene loads a gltf file
func (app *RenderingApp) loadScene(fpath string) error {
	app.sendMessageToClient("loading", fpath)
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

	app.Scene().Add(n)
	root := app.Scene().ChildIndex(n)
	app.nameChildren("/"+strconv.Itoa(root), n)
	app.sendMessageToClient("loaded", fpath)
	return nil
}
