package renderer

import (
	"testing"

	"github.com/g3n/engine/math32"
)

func TestGetCenter(t *testing.T) {
	b := math32.Box3{Min: math32.Vector3{0, 0, 0}, Max: math32.Vector3{10, 10, 100}}
	v := getCenter(b)
	if v.X == 5 && v.Y == 5 && v.Z == 50 {
		return
	} else {
		t.Error("Not Centered")
	}
}

func TestGetViewVectorByName(t *testing.T) {
	v := getViewVectorByName("top")
	if v.Y <= 0 {
		t.Error("top view incorrect")
	}
	v = getViewVectorByName("bottom")
	if v.Y >= 0 {
		t.Error("bottom view incorrect")
	}
	v = getViewVectorByName("left")
	if v.X <= 0 {
		t.Error("left view incorrect")
	}
	v = getViewVectorByName("right")
	if v.X >= 0 {
		t.Error("right view incorrect")
	}
	v = getViewVectorByName("rear")
	if v.Z >= 0 {
		t.Error("rear view incorrect")
	}
	v = getViewVectorByName("front")
	if v.Z <= 0 {
		t.Error("front view incorrect")
	}
}
