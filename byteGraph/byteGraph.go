package byteGraph

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/llgcode/draw2d/draw2dimg"
)

var byteBuffer []int

func AddToByteBuffer(value int) {
	byteBuffer = byteBuffer[1:]
	byteBuffer = append(byteBuffer, value)
}

func convertToRGBA(src *image.NRGBA) *image.RGBA {
	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}

func convertToNRGBA(src *image.RGBA) *image.NRGBA {
	b := src.Bounds()
	img := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	return img
}

func DrawByteGraph(src *image.NRGBA) *image.NRGBA {
	if byteBuffer == nil {
		byteBuffer = make([]int, 30)
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
