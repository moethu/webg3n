package renderer

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

// selectNode uses a raycaster to get the selected node.
// It sends the selection as json to the image channel
// and changes the node's material
func (app *RenderingApp) selectNode(mx float32, my float32, multiselect bool) {
	width, height := app.Window().Size()
	x := (-.5 + mx/float32(width)) * 2.0
	y := (.5 - my/float32(height)) * 2.0
	app.Log().Info("click: %f, %f", x, y)
	r := core.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	app.CameraPersp().SetRaycaster(r, x, y)

	i := r.IntersectObject(app.Scene(), true)

	var object *core.Node
	if len(i) != 0 {
		object = i[0].Object.GetNode()
		app.Log().Info("selected: %s", object.Name())
		app.sendMessageToClient("selected", object.Name())
		if !multiselect {
			app.resetSelection()
		}
		app.changeNodeMaterial(i[0].Object)
	} else {
		if !multiselect {
			app.sendMessageToClient("selected", "")
			app.resetSelection()
		}
	}
}

// resetSelection resets selected nodes to their original state
func (app *RenderingApp) resetSelection() {
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
func (app *RenderingApp) changeNodeMaterial(inode core.INode) {
	gnode, ok := inode.(graphic.IGraphic)

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
