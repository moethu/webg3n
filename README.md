
# g3nserverside

g3nserverside is a GO server side 3D renderer based on the [G3N](https://github.com/g3n/engine) OpenGL engine.
![alt text](https://github.com/moethu/g3nserverside/raw/master/images/screenshot01.png)

## How it works

g3nserverside consists of a GO webserver handling sockets to one G3N rendering instance per connection. The G3N application is constantly streaming images (jpeg) and listens for navigation events and commands from the client. The G3N application is reading [GLTF models](https://github.com/KhronosGroup/glTF), the example model is taken from [here](https://github.com/KhronosGroup/glTF-Sample-Models)

## Why Server Side Rendering

Client and server side renderers have both pros and cons. Depending on your use case server side rendering of 3D models can be a great alternative to client side WebGL rendering. Here are just a few reasons why:

- Browser doesn't need to support WebGL (older Browsers)
- Smaller and weaker devices might not be able to render large models on client side
- The geometry remains on the server and is not transferred to the client

On the other hand it shifts the bottleneck from the client's rendering capabilites to the bandwith.

## Example

This implementation comes with a simple webUI, simply go run it and connect to port 8000.
Once you click connect, a g3n window will appear which is the server side rendering screen.
[Demo video](https://github.com/moethu/g3nserverside/raw/master/images/demo.mp4)

## Contributing

If you find a bug or create a new feature you are encouraged to send pull requests!
