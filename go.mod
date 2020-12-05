//module github.com/moethu/webg3n

module github.com/vshashi01/webg3n

require (
	github.com/disintegration/imaging v1.6.2
	github.com/g3n/engine v0.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/websocket v1.4.2
	github.com/llgcode/draw2d v0.0.0-20200603164053-19660b984a28
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/moethu/webg3n v0.0.0-20201010132405-ce0d591c0e62
	github.com/satori/go.uuid v1.2.0
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace github.com/moethu/webg3n/renderer => ./renderer

replace github.com/moethu/webg3n/byteGraph => ./byteGraph

replace github.com/g3n/engine => github.com/moethu/engine v0.0.0-20200610122637-682e1e061a29

go 1.13
