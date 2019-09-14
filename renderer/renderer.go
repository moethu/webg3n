package renderer

import (
	"log"

	"engine/core"
	"engine/graphic"
	"engine/light"

	"engine/material"
	"engine/math32"
	"engine/util/application"
	"engine/util/logger"
)

type ImageSettings struct {
	jpegQuality int
	saturation  float64
	contrast    float64
	brightness  float64
	blur        float64
	invert      bool
}

type RenderingApp struct {
	application.Application
	x, y, z            float32
	c_imagestream      chan []byte
	c_commands         chan []byte
	Width              int
	Height             int
	imageSettings      ImageSettings
	selectionBuffer    map[core.INode][]graphic.GraphicMaterial
	selection_material material.IMaterial
	modelpath          string
	nodeBuffer         map[string]*core.Node
}

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
		jpegQuality: 60,
		saturation:  0,
		brightness:  0,
		contrast:    0,
		blur:        0,
		invert:      false,
	}

	app.c_imagestream = write
	app.c_commands = read
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
	app.selection_material = material.NewPhong(math32.NewColor("Red"))
	app.selectionBuffer = make(map[core.INode][]graphic.GraphicMaterial)
	app.nodeBuffer = make(map[string]*core.Node)

	app.Gl().ClearColor(1.0, 1.0, 1.0, 1.0)

	er := app.loadScene(app.modelpath)
	if er != nil {
		log.Fatal(er)
	}

	amb := light.NewAmbient(&math32.Color{0.2, 0.2, 0.2}, 1.0)
	app.Scene().Add(amb)

	plight := light.NewPoint(math32.NewColor("white"), 40)
	plight.SetPosition(100, 20, 70)
	plight.SetLinearDecay(.001)
	plight.SetQuadraticDecay(.001)
	app.Scene().Add(plight)

	app.Camera().GetCamera().SetPosition(12, 1, 5)

	p := math32.Vector3{X: 0, Y: 0, Z: 0}
	app.Camera().GetCamera().LookAt(&p)
	app.CameraPersp().SetFov(65)
	app.zoomToExtent()
	app.Orbit().Enabled = true
	app.Application.Subscribe(application.OnAfterRender, app.onRender)
}
