package renderer

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
)

// byte buffer holding sent bytes
var byteBuffer []int

// AddToByteBuffer adding sent bytes to the history stack
func AddToByteBuffer(value int) {
	byteBuffer = byteBuffer[1:]
	byteBuffer = append(byteBuffer, value)
}

// DrawByteGraph draws the histroy of sent bytes onto the image
func DrawByteGraph(img *image.RGBA) *image.RGBA {
	if byteBuffer == nil {
		byteBuffer = make([]int, 50)
	}
	gc := draw2dimg.NewGraphicContext(img)
	x := float64(len(byteBuffer)) * 3.0
	y := float64(img.Bounds().Dy())
	gc.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	gc.SetLineWidth(0.5)

	// draw 100 kb line
	gc.BeginPath()
	gc.MoveTo(0, y-100)
	gc.LineTo(x, y-100)
	gc.Close()
	gc.Stroke()

	gc.SetStrokeColor(color.Black)
	gc.SetLineWidth(0.5)
	// draw sent bytes as bars in kb
	for _, value := range byteBuffer {
		gc.BeginPath()
		gc.MoveTo(x, y)
		gc.LineTo(x, y-float64(value/1000))
		gc.Close()
		gc.Stroke()
		x = x - 3.0
	}
	return img
}
