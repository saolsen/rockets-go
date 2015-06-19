package main

// @TODO: Make sure to statically link nanovg
// @TODO: Make sure the font file gets loaded in some way that works for releasing.

// #cgo CFLAGS: -Inanovg/src
// #cgo LDFLAGS: -framework OpenGL -Lnanovg/build -lnanovg
// #include <OpenGL/gl3.h>
// #include "nanovg.h"
// #define NANOVG_GL3_IMPLEMENTATION
// #include "nanovg_gl.h"
import "C"

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const WindowWidth = 1280
const WindowHeight = 720

func f64round(f float64) float64 {
	return math.Floor(f + 0.5)
}

func f64round2(f float64) float64 {
	return f64round(f*100) / 100
}

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Rockets", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Init Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Opengl version: ", version)

	// Setup nvg
	vg := C.nvgCreateGL3(C.NVG_ANTIALIAS | C.NVG_STENCIL_STROKES | C.NVG_DEBUG)
	if vg == nil {
		panic("Could not init nanovg.")
	}

	C.nvgCreateFont(vg,
		C.CString("basic"),
		C.CString("SourceSansPro-Regular.ttf"))

	// setup game
	state := Setup()

	var lastStartTime, startTime, frameTime, gameCodeTime float64
	lastStartTime = glfw.GetTime()

	for !window.ShouldClose() {
		startTime = glfw.GetTime()
		frameTime = startTime - lastStartTime
		lastStartTime = startTime

		// Do Opengl stuff
		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		C.nvgBeginFrame(vg, WindowWidth, WindowHeight, 2.0)

		UpdateAndRender(vg, state, frameTime)

		gameCodeTime = glfw.GetTime() - startTime

		// Display Stats
		stats := fmt.Sprintf("game: %.2fms, frame: %.2fms, fps: %.2f",
			f64round2(gameCodeTime*100.0), f64round2(frameTime*100.0), f64round2(1/frameTime))

		C.nvgSave(vg)
		C.nvgFontSize(vg, 14)
		C.nvgFillColor(vg, C.nvgRGBf(1, 1, 1))
		C.nvgText(vg, 5, WindowHeight-10, C.CString(stats), nil)
		C.nvgRestore(vg)

		C.nvgEndFrame(vg)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
