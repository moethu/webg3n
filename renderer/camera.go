package renderer

import (
	"log"

	"github.com/g3n/engine/math32"
)

type Standardview int

// Standard Views
const (
	Front  = Standardview(0)
	Bottom = Standardview(1)
	Left   = Standardview(2)
	Right  = Standardview(3)
	Rear   = Standardview(4)
	Top    = Standardview(5)
)

// getCenter returns boundingbox centerpoint
func getCenter(box math32.Box3) *math32.Vector3 {
	return box.Center(nil)
}

// FocusOnSelection focuses camera on currently selected elements
func (app *RenderingApp) FocusOnSelection() {
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

	position := app.camera.Position()
	app.extentViewOn(bbox, position)
}

// SetStandardView sets the camera view to a standard
func (app *RenderingApp) SetStandardView(view Standardview) {
	modifier := math32.Vector3{X: 0, Y: 0, Z: 0}
	switch view {
	case Top:
		modifier.Y = 10
	case Bottom:
		modifier.Y = -10
	case Front:
		modifier.Z = 10
	case Rear:
		modifier.Z = -10
	case Left:
		modifier.X = 10
	case Right:
		modifier.X = -10
	default:
		return
	}
	bbox := app.scene.ChildAt(0).BoundingBox()
	C := bbox.Center(nil)
	pos := modifier.Add(C)
	app.extentViewOn(&bbox, *pos)
}

// extentViewOn extents the camera view on a bounding box searching for the best fit
// camera position between the centerpoint and the position
func (app *RenderingApp) extentViewOn(bbox *math32.Box3, position math32.Vector3) {
	C := bbox.Center(nil)
	r := C.DistanceTo(&bbox.Max)
	a := app.camera.Fov()
	d := r / math32.Sin(a/2)
	P := math32.Vector3{X: C.X, Y: C.Y, Z: C.Z}
	dir := math32.Vector3{X: C.X, Y: C.Y, Z: C.Z}
	P.Add(((position.Sub(C)).Normalize().MultiplyScalar(d * -1)))
	dir.Sub(&P)
	app.camera.SetPositionVec(&P)
	up := math32.Vector3{X: 0, Y: 1, Z: 0}
	app.camera.LookAt(C, &up)
}

// ZoomExtent zooms to sceene extents
func (app *RenderingApp) ZoomExtent() {
	pos := app.camera.Position()
	bbox := app.scene.ChildAt(0).BoundingBox()
	app.extentViewOn(&bbox, pos)
}
