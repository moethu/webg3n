package renderer

import (
	"testing"

	"github.com/g3n/engine/window"
)

func assert(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Error("Actual:", actual, "Expected:", expected)
	}
}
func TestMapMouseButton(t *testing.T) {
	assert(t, mapMouseButton("0"), window.MouseButtonLeft)
	assert(t, mapMouseButton("1"), window.MouseButtonMiddle)
	assert(t, mapMouseButton("2"), window.MouseButtonRight)
	assert(t, mapMouseButton("3"), window.MouseButtonLeft)
}

func TestMapKey(t *testing.T) {
	assert(t, mapKey("38"), window.KeyUp)
	assert(t, mapKey("37"), window.KeyLeft)
	assert(t, mapKey("39"), window.KeyRight)
	assert(t, mapKey("40"), window.KeyDown)
	assert(t, mapKey("41"), window.KeyEnter)
}

func TestGetValueInRange(t *testing.T) {
	assert(t, getValueInRange(6, 1, 5), 5)
	assert(t, getValueInRange(3, 1, 5), 3)
	assert(t, getValueInRange(0, 1, 5), 1)
}
