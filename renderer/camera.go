package renderer

import "engine/math32"

func getCenter(box math32.Box3) *math32.Vector3 {
	return box.Center(nil)
}

func (app *RenderingApp) focusOnElement() {
	for inode, _ := range app.selectionBuffer {
		pos := getCenter(inode.BoundingBox())
		app.Camera().GetCamera().LookAt(pos)
	}
}

func (app *RenderingApp) setCameraTop() {
	pos := app.Camera().GetCamera().Position()
	pos.Y -= 8
	app.Camera().GetCamera().LookAt(&pos)
}
