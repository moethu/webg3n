package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"time"

	"engine/util/application"
	"engine/util/logger"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/webg3n", serveWebsocket)
	http.HandleFunc("/", home)
	log.Println("Starting HTTP Server on Port 8000")
	http.ListenAndServe(*addr, nil)
}

var addr = flag.String("addr", "0.0.0.0:8000", "http service address")
var upgrader = websocket.Upgrader{}

type RData struct {
	Image string
	Stamp time.Time
}

func load3DApplication(app *renderingApp, h int, w int, write chan []byte, read chan []byte, modelpath string) {
	a, err := application.Create(application.Options{
		Title:       "ServerSideTutorial01",
		Width:       w,
		Height:      h,
		Fullscreen:  false,
		LogPrefix:   "ServerSide01",
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
	app.c_imagestream = write
	app.c_commands = read
	app.jpegQuality = 60
	app.modelpath = modelpath
	app.setupScene()
	go app.commandLoop()
	err = app.Run()
	if err != nil {
		panic(err)
	}

	app.Log().Info("app was running for %f seconds\n", app.RunSeconds())
}

// Home route, loading template and serving it
func home(w http.ResponseWriter, r *http.Request) {
	viewertemplate := template.Must(template.ParseFiles("templates/viewer.html"))
	viewertemplate.Execute(w, "ws://"+r.Host+"/webg3n")
}
