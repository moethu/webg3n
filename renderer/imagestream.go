package renderer

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

// onRender event handler for onRender event
func (a *RenderingApp) onRender(evname string, ev interface{}) {
	a.makeScreenShot()
}

// makeScreenShot reads the opengl buffer, encodes it as jpeg and sends it to the channel
func (app *RenderingApp) makeScreenShot() {
	w := app.Width
	h := app.Height
	data := app.Gl().ReadPixels(0, 0, w, h, 6408, 5121)
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	img.Pix = data

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
	buf := new(bytes.Buffer)
	var opt jpeg.Options
	opt.Quality = app.imageSettings.jpegQuality
	jpeg.Encode(buf, img, &opt)
	imageBit := buf.Bytes()

	// could use encode directly
	imgBase64Str := base64.StdEncoding.EncodeToString([]byte(imageBit))
	app.c_imagestream <- []byte(imgBase64Str)
}
