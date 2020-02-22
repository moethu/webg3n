module github.com/moethu/webg3n

require (
	github.com/disintegration/imaging v1.6.1
	github.com/g3n/engine v0.1.0
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/websocket v1.4.1
	github.com/llgcode/draw2d v0.0.0-20200110163050-b96d8208fcfc
	github.com/satori/go.uuid v1.2.0
)

replace github.com/moethu/webg3n/renderer => ./renderer

replace github.com/moethu/webg3n/byteGraph => ./byteGraph

replace github.com/g3n/engine => github.com/moethu/engine v0.0.0-20190918211458-57b17b524856

go 1.13
