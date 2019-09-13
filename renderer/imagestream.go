package renderer

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

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

	if app.brightness != 0 {
		img = imaging.AdjustBrightness(img, app.brightness)
	}
	if app.contrast != 0 {
		img = imaging.AdjustContrast(img, app.contrast)
	}
	if app.saturation != 0 {
		img = imaging.AdjustSaturation(img, app.saturation)
	}
	if app.blur != 0 {
		img = imaging.Blur(img, app.blur)
	}
	if app.invert {
		img = imaging.Invert(img)
	}

	img = imaging.FlipV(img)
	buf := new(bytes.Buffer)
	var opt jpeg.Options
	opt.Quality = app.jpegQuality
	jpeg.Encode(buf, img, &opt)
	imageBit := buf.Bytes()
	imgBase64Str := base64.StdEncoding.EncodeToString([]byte(imageBit))
	app.c_imagestream <- []byte(imgBase64Str)
}
