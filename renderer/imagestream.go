package renderer

import (
	"image"
	"log"

	"github.com/moethu/imaging"
	"github.com/moethu/webg3n/encoders"
)

// onRender event handler for onRender event
func (app *RenderingApp) onRender(evname string, ev interface{}) {
	app.makeScreenShot()
}

var md5SumBuffer [16]byte
var es encoders.Service
var e encoders.Encoder

// makeScreenShot reads the opengl buffer, encodes it as jpeg and sends it to the channel
func (app *RenderingApp) makeScreenShot() {
	w := app.Width
	h := app.Height
	data := app.Gl().ReadPixels(0, 0, w, h, 6408, 5121)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Pix = data

	if app.imageSettings.getPixelation() > 1.0 {
		img = imaging.Fit(img, int(float64(w)/app.imageSettings.getPixelation()), int(float64(h)/app.imageSettings.getPixelation()), imaging.NearestNeighbor)
	}
	if app.imageSettings.brightness != 0 {
		img = imaging.AdjustBrightness(img, app.imageSettings.brightness)
	}
	if app.imageSettings.contrast != 0 {
		img = imaging.AdjustContrast(img, app.imageSettings.contrast)
	}
	if app.imageSettings.saturation != 0 {
		img = imaging.AdjustSaturation(img, app.imageSettings.saturation)
	}
	if app.imageSettings.blur != 0 {
		img = imaging.Blur(img, app.imageSettings.blur)
	}
	if app.imageSettings.invert {
		img = imaging.Invert(img)
	}

	img = imaging.FlipV(img)

	if app.Debug {
		img = DrawByteGraph(img)
	}

	if es == nil {
		es = encoders.NewEncoderService()
		e, _ = es.NewEncoder(encoders.VP8Codec, image.Point{X: w, Y: h}, 20)
	}

	d, err := e.Encode(img)
	if err != nil {
		log.Println(err)
	}

	app.cImagestream <- d
}
