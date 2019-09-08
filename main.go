package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
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

// Home route, loading template and serving it
func home(w http.ResponseWriter, r *http.Request) {
	viewertemplate := template.Must(template.ParseFiles("templates/viewer.html"))
	viewertemplate.Execute(w, r.Host)
}
