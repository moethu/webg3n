
# webg3n

webg3n is a 3D web-viewer running [G3N](https://github.com/g3n/engine) as a server side OpenGL renderer.

[![Go Report Card](https://goreportcard.com/badge/github.com/moethu/webg3n)](https://goreportcard.com/report/github.com/moethu/webg3n)

![alt text](https://github.com/moethu/g3nserverside/raw/master/images/screenshot01.png)

[Checkout this demo video](https://vimeo.com/358812535)

## How it works

![alt text](https://github.com/moethu/g3nserverside/raw/master/images/arc.png)

webg3n is a GO webserver handling sockets to one G3N rendering instance per connection. The G3N application is constantly streaming images (jpeg) and listens for navigation events and commands from the client. The server supports multiple clients at the same time. For each client it will open one socket and spin off a separate 3D engine for rendering. The G3N application is reading [GLTF models](https://github.com/KhronosGroup/glTF), the example model is taken from [here](https://github.com/KhronosGroup/glTF-Sample-Models)

## Why Server Side Rendering

Client and server side renderers have both pros and cons. Depending on your use case server side rendering of 3D models can be a great alternative to client side WebGL rendering. Here are just a few reasons why:

- Browser doesn't need to support WebGL (older Browsers)
- Smaller and weaker devices might not be able to render large models on client side
- The geometry remains on the server and is not transferred to the client

On the other hand it shifts the bottleneck from the client's rendering capabilites to the bandwith.

## Dependencies

Go 1.8+ is required. The engine also requires the system to have an OpenGL driver and a GCC-compatible C compiler.

Requires this modified [G3N engine](https://github.com/moethu/engine) which gets installed via go modules.

On Unix-based systems the engine depends on some C libraries that can be installed using the appropriate distribution package manager. See below for OS specific requirements.

#### Ubuntu/Debian-like

```
$ sudo apt-get install xorg-dev libgl1-mesa-dev libopenal1 libopenal-dev libvorbis0a libvorbis-dev libvorbisfile3
```

#### Fedora

```
$ sudo dnf -y install xorg-x11-proto-devel mesa-libGL mesa-libGL-devel openal-soft openal-soft-devel libvorbis libvorbis-devel glfw-devel libXi-devel
```

#### CentOS 7

Enable the EPEL repository:
```
$ sudo yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
```
Then install the same packages as for Fedora - remember to use yum instead of dnf for the package installation command.

#### Windows

The necessary audio libraries sources and DLLs are supplied but they need to be installed manually. Please see Audio libraries for Windows for details. We tested the Windows build using the mingw-w64 toolchain (you can download this file in particular).

#### macOS

Install the development files of OpenAL and Vorbis using Homebrew:
```
brew install libvorbis openal-soft
```

## Example

This implementation comes with a simple webUI, simply go run it and connect to port 8000.
Once you click connect, a g3n window will appear which is the server side rendering screen.

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
