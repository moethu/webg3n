package renderer

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/loader/gltf"
)

// nameChildren names all gltf nodes by path
func (app *RenderingApp) nameChildren(p string, n core.INode) {
	node := n.GetNode()
	node.SetName(p)
	//app.nodeBuffer[p] = node
	for _, child := range node.Children() {
		idx := node.ChildIndex(child)
		title := p + "/" + strconv.Itoa(idx)
		app.nameChildren(title, child)
	}
}

// LoadScene loads a gltf file
func (app *RenderingApp) LoadScene(fpath string) (*gltf.GLTF, error) {
	app.SendMessageToClient("loading", fpath)
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
		return nil, fmt.Errorf("unrecognized file extension:%s", ext)
	}

	if err != nil {
		return nil, err
	}

	defaultSceneIdx := 0
	if g.Scene != nil {
		defaultSceneIdx = *g.Scene
	}

	// Create default scene
	n, err := g.LoadScene(defaultSceneIdx)
	if err != nil {
		return nil, err
	}

	//meshList := make([]*graphic.Mesh, 10)

	mesh := returnFirstMesh(n)

	if mesh != nil {
		slc := strings.Split(fpath, ".")
		ok := app.LoadMeshEntity(mesh, slc[0])

		if ok {
			app.zoomToExtent()
			app.SendMessageToClient("loaded", fpath)
		}

	}

	//n.GetNode().SetName(fpath)
	//n.GetNode().GetNode().SetName(fpath)
	//app.Scene().Add(n)
	//root := app.Scene().ChildIndex(n)
	//app.nameChildren("/"+strconv.Itoa(root), n)

	return g, nil
}

func returnFirstMesh(node core.INode) *graphic.Mesh {
	mesh, ok := node.(*graphic.Mesh)

	if ok {
		return mesh
	}

	for _, ci := range node.GetNode().Children() {
		return returnFirstMesh(ci)

	}
	return nil
}
