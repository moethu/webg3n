module github.com/moethu/webg3n

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/g3n/engine v0.1.0
	github.com/gin-gonic/gin v1.7.0
	github.com/golang/snappy v0.0.2 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/llgcode/draw2d v0.0.0-20200603164053-19660b984a28
	github.com/mholt/archiver v3.1.1+incompatible // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/moethu/imaging v1.6.3
	github.com/nickalie/go-binwrapper v0.0.0-20190114141239-525121d43c84 // indirect
	github.com/nickalie/go-mozjpegbin v0.0.0-20170427050522-d8a58e243a3d
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pixiv/go-libjpeg v0.0.0-20190822045933-3da21a74767d
	github.com/satori/go.uuid v1.2.0
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/net v0.7.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace github.com/moethu/webg3n/renderer => ./renderer

replace github.com/moethu/webg3n/byteGraph => ./byteGraph

replace github.com/g3n/engine => github.com/moethu/engine v0.0.0-20200610122637-682e1e061a29

go 1.13
