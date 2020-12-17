package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/g3n/engine/window"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	port := ":8000"
	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	//router.Static("/static/", "./static/")
	router.Any("/webg3n", serveWebsocket)
	router.PUT("/loadModel", loadModel)
	router.GET("/objects", getObjects)
	//router.GET("/", home)
	log.Println("Starting HTTP Server on Port 8000")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	window.DestroyGlfwManager()
	log.Println("Server exiting")
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

// Home route, loading template and serving it
func home(c *gin.Context) {
	viewertemplate := template.Must(template.ParseFiles("templates/webg3n.html"))
	viewertemplate.Execute(c.Writer, c.Request.Host)
}
