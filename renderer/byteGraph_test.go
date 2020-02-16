package renderer

import (
	"testing"
)

func TestSlice(t *testing.T) {
	byteBuffer = []int{1, 2, 3, 4, 5, 6}
	AddToByteBuffer(7)
	assert(t, byteBuffer[5], 7)
	assert(t, len(byteBuffer), 6)
	AddToByteBuffer(8)
	assert(t, byteBuffer[5], 8)
	assert(t, len(byteBuffer), 6)
	AddToByteBuffer(9)
	assert(t, byteBuffer[5], 9)
	assert(t, len(byteBuffer), 6)
}
