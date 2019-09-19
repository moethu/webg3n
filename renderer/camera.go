package renderer

import (
	"github.com/g3n/engine/math32"
	"log"
)

func getCenter(box math32.Box3) *math32.Vector3 {
	return box.Center(nil)
}

func (app *RenderingApp) focusOnSelection() {
	var bbox *math32.Box3
	first := true
	for inode, _ := range app.selectionBuffer {
		tmp := inode.BoundingBox()
		if first {
			bbox = math32.NewBox3(&tmp.Min, &tmp.Max)
			log.Println(bbox)
			first = false
		} else {
			bbox.ExpandByPoint(&tmp.Min)
			bbox.ExpandByPoint(&tmp.Max)
		}
	}
	if first {
		return
	}
	position := app.Camera().GetCamera().Position()
	C := bbox.Center(nil)
	r := C.DistanceTo(&bbox.Max)
	a := app.CameraPersp().Fov()
	d := r / math32.Sin(a/2)
	P := math32.Vector3{X: C.X, Y: C.Y, Z: C.Z}
	dir := math32.Vector3{X: C.X, Y: C.Y, Z: C.Z}
	P.Add(((position.Sub(C)).Normalize().MultiplyScalar(d)))
	dir.Sub(&P)
	app.Camera().GetCamera().SetPositionVec(&P)
	app.Camera().GetCamera().LookAt(C)
}

func (app *RenderingApp) setCamera(view string) {
	modifier := math32.Vector3{X: 0, Y: 0, Z: 0}
	switch view {
	case "top":
		modifier.Y = 10
	case "bottom":
		modifier.Y = -10
	case "front":
		modifier.Z = 10
	case "rear":
		modifier.Z = -10
	case "left":
		modifier.X = 10
	case "right":
		modifier.X = -10
	}
	bbox := app.Scene().ChildAt(0).BoundingBox()
	C := bbox.Center(nil)
	pos := modifier.Add(C)
	app.focusCameraToCenter(*pos)
}

func (app *RenderingApp) focusCameraToCenter(position math32.Vector3) {
	bbox := app.Scene().ChildAt(0).BoundingBox()
	C := bbox.Center(nil)
	r := C.DistanceTo(&bbox.Max)
	a := app.CameraPersp().Fov()
	d := r / math32.Sin(a/2)
	P := math32.Vector3{X: C.X, Y: C.Y, Z: C.Z}
	dir := math32.Vector3{X: C.X, Y: C.Y, Z: C.Z}
	P.Add(((position.Sub(C)).Normalize().MultiplyScalar(d)))
	dir.Sub(&P)
	app.Camera().GetCamera().SetPositionVec(&P)
	app.Camera().GetCamera().LookAt(C)
}

func (app *RenderingApp) zoomToExtent() {
	pos := app.Camera().GetCamera().Position()
	app.focusCameraToCenter(pos)
}
