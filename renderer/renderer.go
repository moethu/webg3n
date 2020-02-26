package renderer

import (
	"log"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	"github.com/g3n/engine/util/logger"
)

// ImageSettings for rendering image
type ImageSettings struct {
	saturation   float64
	contrast     float64
	brightness   float64
	blur         float64
	pixelation   float64
	invert       bool
	quality      Quality
	isNavigating bool
}

// getJpegQuality returns quality depending on navigation movement
func (i *ImageSettings) getJpegQuality() int {
	if i.isNavigating {
		return i.quality.jpegQualityNav
	} else {
		return i.quality.jpegQualityStill
	}
}

// getPixelation returns pixelation depending on navigation movement
// A global pixelation level will override preset pixelation levels
func (i *ImageSettings) getPixelation() float64 {
	if i.pixelation > 1.0 {
		return i.pixelation
	}
	if i.isNavigating {
		return i.quality.pixelationNav
	} else {
		return i.quality.pixelationStill
	}
}

// Quality Image quality settings for still and navigating situations
type Quality struct {
	jpegQualityStill int
	jpegQualityNav   int
	pixelationStill  float64
	pixelationNav    float64
}

// high image quality definition
var highQ Quality = Quality{jpegQualityStill: 100, jpegQualityNav: 90, pixelationStill: 1.0, pixelationNav: 1.0}

// medium image quality definition
var mediumQ Quality = Quality{jpegQualityStill: 80, jpegQualityNav: 60, pixelationStill: 1.0, pixelationNav: 1.2}

// low image quality definition
var lowQ Quality = Quality{jpegQualityStill: 60, jpegQualityNav: 40, pixelationStill: 1.0, pixelationNav: 1.5}

// RenderingApp application settings
type RenderingApp struct {
	application.Application
	x, y, z           float32
	cImagestream      chan []byte
	cCommands         chan []byte
	Width             int
	Height            int
	imageSettings     ImageSettings
	selectionBuffer   map[core.INode][]graphic.GraphicMaterial
	selectionMaterial material.IMaterial
	modelpath         string
	nodeBuffer        map[string]*core.Node
	Debug             bool
}

// LoadRenderingApp loads the rendering application
func LoadRenderingApp(app *RenderingApp, sessionId string, h int, w int, write chan []byte, read chan []byte, modelpath string) {
	a, err := application.Create(application.Options{
		Title:       "g3nServerApplication",
		Width:       w,
		Height:      h,
		Fullscreen:  false,
		LogPrefix:   sessionId,
		LogLevel:    logger.DEBUG,
		TargetFPS:   30,
		EnableFlags: true,
	})

	if err != nil {
		panic(err)
	}

	app.Application = *a
	app.Width = w
	app.Height = h

	app.imageSettings = ImageSettings{
		saturation: 0,
		brightness: 0,
		contrast:   0,
		blur:       0,
		pixelation: 1.0,
		invert:     false,
		quality:    highQ,
	}

	app.cImagestream = write
	app.cCommands = read
	app.modelpath = modelpath
	app.setupScene()
	go app.commandLoop()
	err = app.Run()
	if err != nil {
		panic(err)
	}

	app.Log().Info("app was running for %f seconds\n", app.RunSeconds())
}

// setupScene sets up the current scene
func (app *RenderingApp) setupScene() {
	app.selectionMaterial = material.NewPhong(math32.NewColor("Red"))
	app.selectionBuffer = make(map[core.INode][]graphic.GraphicMaterial)
	app.nodeBuffer = make(map[string]*core.Node)

	app.Gl().ClearColor(1.0, 1.0, 1.0, 1.0)

	er := app.loadScene(app.modelpath)
	if er != nil {
		log.Fatal(er)
	}

	amb := light.NewAmbient(&math32.Color{R: 0.2, G: 0.2, B: 0.2}, 1.0)
	app.Scene().Add(amb)

	plight := light.NewPoint(math32.NewColor("white"), 40)
	plight.SetPosition(100, 20, 70)
	plight.SetLinearDecay(.001)
	plight.SetQuadraticDecay(.001)
	app.Scene().Add(plight)

	app.Camera().GetCamera().SetPosition(12, 1, 5)

	p := math32.Vector3{X: 0, Y: 0, Z: 0}
	app.Camera().GetCamera().LookAt(&p)
	app.CameraPersp().SetFov(50)
	app.zoomToExtent()
	app.Orbit().Enabled = true
	app.Application.Subscribe(application.OnAfterRender, app.onRender)
}
