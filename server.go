package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"webg3n/renderer"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	writeTimeout   = 10 * time.Second
	readTimeout    = 60 * time.Second
	pingPeriod     = (readTimeout * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	app renderer.RenderingApp

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channels messages.
	write chan []byte
	read  chan []byte
}

// streamReader reads messages from the websocket connection and fowards them to the read channel
func (c *Client) streamReader() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readTimeout))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(readTimeout)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.read <- message
	}
}

// streamWriter writes messages from the write channel to the websocket connection
func (c *Client) streamWriter() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.write:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.write)
			for i := 0; i < n; i++ {
				w.Write(<-c.write)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWebsocket handles websocket requests from the peer.
func serveWebsocket(c *gin.Context) {
	sessionId := uuid.NewV4()
	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn.EnableWriteCompression(true)

	// create two channels for read write concurrency
	c_write := make(chan []byte)
	c_read := make(chan []byte)

	client := &Client{conn: conn, write: c_write, read: c_read}

	// get scene width and height from url query params
	// default to 800 if they are not set
	height, err := strconv.Atoi(c.Request.URL.Query().Get("h"))
	if err != nil {
		log.Println(err)
		height = 800
	}
	width, err := strconv.Atoi(c.Request.URL.Query().Get("w"))
	if err != nil {
		log.Println(err)
		width = 800
	}

	modelPath := "models/"
	model := c.Request.URL.Query().Get("model")
	if model == "" {
		model = "Cathedral.glb"
	}
	if _, err := os.Stat(modelPath + model); os.IsNotExist(err) {
		model = "Cathedral.glb"
	}

	// run 3d application in separate go routine
	// this is currently not threadafe but it's a single 3d app per socket
	go renderer.LoadRenderingApp(&client.app, sessionId.String(), height, width, c_write, c_read, modelPath+model)

	// run reader and writer in two different go routines
	// so they can act concurrently
	go client.streamReader()
	go client.streamWriter()
}
