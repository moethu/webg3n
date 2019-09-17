package renderer

import (
	"log"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/logger"
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
	*app.Application
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
	log                *logger.Logger // Application logger
	scene              *core.Node     // Scene rendered

	camera *camera.Camera // Current camera
	orbit  *OrbitControl  // Camera orbit controller
}

func LoadRenderingApp(sessionId string, h int, w int, write chan []byte, read chan []byte, modelpath string) {

	a := new(RenderingApp)
	a.Application = app.App()

	a.Width = w
	a.Height = h

	a.imageSettings = ImageSettings{
		jpegQuality: 60,
		saturation:  0,
		brightness:  0,
		contrast:    0,
		blur:        0,
		invert:      false,
	}

	a.c_imagestream = write
	a.c_commands = read
	a.modelpath = modelpath
	a.SetupScene()
	go a.commandLoop()
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(a.scene, a.camera)
		a.makeScreenShot()
	})
	a.log.Info("App closed\n")
}

// setupScene sets up the current scene
func (app *RenderingApp) SetupScene() {

	// Creates application logger
	app.log = logger.New("G3ND", nil)
	app.log.AddWriter(logger.NewConsole(false))
	app.log.SetFormat(logger.FTIME | logger.FMICROS)
	app.log.SetLevel(logger.DEBUG)
	app.log.Info("Starting")

	// Define Selection Material and buffer
	app.selection_material = material.NewPhong(math32.NewColor("Red"))
	app.selectionBuffer = make(map[core.INode][]graphic.GraphicMaterial)
	app.nodeBuffer = make(map[string]*core.Node)

	// setup new scene and white background
	app.scene = core.NewNode()
	app.Gls().ClearColor(1.0, 1.0, 1.0, 1.0)

	// load model from file
	er := app.loadScene(app.modelpath)
	if er != nil {
		log.Fatal(er)
	}

	// Create perspective camera
	width, height := app.GetSize()
	aspect := float32(width) / float32(height)
	app.camera = camera.New(aspect)
	app.scene.Add(app.camera) // Add camera to scene (important for audio demos)

	amb := light.NewAmbient(&math32.Color{0.2, 0.2, 0.2}, 1.0)
	app.scene.Add(amb)

	plight := light.NewPoint(math32.NewColor("white"), 40)
	plight.SetPosition(100, 20, 70)
	plight.SetLinearDecay(.001)
	plight.SetQuadraticDecay(.001)
	app.scene.Add(plight)

	app.ZoomExtent()
	app.orbit = NewOrbitControl(app.camera)
}
