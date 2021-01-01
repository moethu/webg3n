package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/moethu/webg3n/renderer"

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

// Client holding g3napp, socket and channels
type Client struct {
	app renderer.RenderingApp

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channels messages.
	write chan []byte // images and data to client
	read  chan []byte // commands from client
}

// streamReader reads messages from the websocket connection and fowards them to the read channel
func (c *Client) streamReader() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readTimeout))
	// SetPongHandler sets the handler for pong messages received from the peer.
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(readTimeout)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// feed message to command channel
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
		// Go’s select lets you wait on multiple channel operations.
		// We’ll use select to await both of these values simultaneously.
		select {
		case message, ok := <-c.write:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// NextWriter returns a writer for the next message to send.
			// The writer's Close method flushes the complete message to the network.
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.write)
			for i := 0; i < n; i++ {
				w.Write(<-c.write)
			}

			if err := w.Close(); err != nil {
				return
			}

		//a channel that will send the time with a period specified by the duration argument
		case <-ticker.C:
			// SetWriteDeadline sets the deadline for future Write calls
			// and any currently-blocked Write call.
			// Even if write times out, it may return n > 0, indicating that
			// some of the data was successfully written.
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
	cWrite := make(chan []byte)
	cRead := make(chan []byte)

	client := &Client{conn: conn, write: cWrite, read: cRead}

	// get scene width and height from url query params
	// default to 800 if they are not set
	height := getParameterDefault(c, "h", 800)
	width := getParameterDefault(c, "w", 800)

	modelPath := "models/"
	defaultModel := "Cathedral.glb"
	model := c.Request.URL.Query().Get("model")
	if model == "" {
		model = defaultModel
	}
	if _, err := os.Stat(modelPath + model); os.IsNotExist(err) {
		model = defaultModel
	}

	// run 3d application in separate go routine
	go renderer.LoadRenderingApp(&client.app, sessionId.String(), height, width, cWrite, cRead, modelPath+model)

	// run reader and writer in two different go routines
	// so they can act concurrently
	go client.streamReader()
	go client.streamWriter()
}

// getParameterDefault gets a parameter and returns default value if its not set
func getParameterDefault(c *gin.Context, name string, defaultValue int) int {
	val, err := strconv.Atoi(c.Request.URL.Query().Get(name))
	if err != nil {
		log.Println(err)
		return defaultValue
	}
	return val
}
