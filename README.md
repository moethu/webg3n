
# webg3n

webg3n is a 3D web-viewer running [G3N](https://github.com/g3n/engine) as a server side OpenGL renderer.

[![Go Report Card](https://goreportcard.com/badge/github.com/moethu/webg3n)](https://goreportcard.com/report/github.com/moethu/webg3n)

![alt text](https://github.com/moethu/g3nserverside/raw/master/images/screenshot03.png)

## How it works

![alt text](https://github.com/moethu/g3nserverside/raw/master/images/arc.png)

webg3n is a GO webserver handling sockets to one G3N rendering instance per connection. The G3N application is constantly streaming images (jpeg) and listens for navigation events and commands from the client. The server supports multiple clients at the same time. For each client it will open one socket and spin off a separate 3D engine for rendering. The G3N application is reading [GLTF models](https://github.com/KhronosGroup/glTF), the example model is taken from [here](https://github.com/KhronosGroup/glTF-Sample-Models)

## Why Server Side Rendering

Client and server side renderers have both pros and cons. Depending on your use case server side rendering of 3D models can be a great alternative to client side WebGL rendering. Here are just a few reasons why:

- Browser doesn't need to support WebGL (older Browsers)
- Smaller and weaker devices might not be able to render large models on client side
- The geometry remains on the server and is not transferred to the client

On the other hand it shifts the bottleneck from the client's rendering capabilites to the bandwith.

## Requirements

Requires this modified [G3N engine](https://github.com/moethu/engine) and its [dependencies](https://github.com/moethu/engine#dependencies)

## Example

This implementation comes with a simple webUI, simply go run it and connect to port 8000.
Once you click connect, a g3n window will appear which is the server side rendering screen.
[Demo video](https://vimeo.com/358812535)

## Features

- mouse navigation
- keyboard navigation
- multi-select elements
- hide and unhide element selection
- focus on element selection
- zoom extents
- default views (top, down, left, right, front, rear)
- change field of view
- set compression quality
- automatically set higher compression while navigating
- adjust image settings (invert, brightness, contrast, saturation, blur)

## Contributing

If you find a bug or create a new feature you are encouraged to send pull requests!
