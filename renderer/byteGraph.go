package renderer

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/llgcode/draw2d/draw2dimg"
)

// byte buffer holding sent bytes
var byteBuffer []int

// AddToByteBuffer adding sent bytes to the history stack
func AddToByteBuffer(value int) {
	byteBuffer = byteBuffer[1:]
	byteBuffer = append(byteBuffer, value)
}

// convertToRGBA converts NRGBA to RGBA
func convertToRGBA(src *image.NRGBA) *image.RGBA {
	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}

// convertToNRGBA converts RGBA to NRGBA
func convertToNRGBA(src *image.RGBA) *image.NRGBA {
	b := src.Bounds()
	img := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}

// DrawByteGraph draws the histroy of sent bytes onto the image
func DrawByteGraph(src *image.NRGBA) *image.NRGBA {
	if byteBuffer == nil {
		byteBuffer = make([]int, 50)
	}
	img := convertToRGBA(src)
	gc := draw2dimg.NewGraphicContext(img)
	x := float64(len(byteBuffer)) * 3.0
	y := float64(img.Bounds().Dy())
	gc.SetStrokeColor(color.Black)

	// draw 100 kb line
	gc.BeginPath()
	gc.MoveTo(0, y-100)
	gc.LineTo(x, y-100)
	gc.Close()
	gc.Stroke()

	// draw sent bytes as bars in kb
	for _, value := range byteBuffer {
		gc.BeginPath()
		gc.MoveTo(x, y)
		gc.LineTo(x, y-float64(value/1000))
		gc.Close()
		gc.Stroke()
		x = x - 3.0
	}
	return convertToNRGBA(img)
}
