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

func f64round(f float64) float64 {
	return math.Floor(f + 0.5)
}

func f64round2(f float64) float64 {
	return f64round(f*100) / 100
}

const WindowWidth = 1280
const WindowHeight = 720

type InputState struct {
	last_mouse_x    float64
	last_mouse_y    float64
	current_mouse_x float64
	current_mouse_y float64
	click           bool
	start_dragging  bool
	is_dragging     bool
	end_dragging    bool
}

type GuiState struct {
	vg    *C.struct_NVGcontext
	input *InputState
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
	input := InputState{}
	gui := GuiState{vg, &input}

	var lastStartTime, startTime, frameTime, gameCodeTime float64
	lastStartTime = glfw.GetTime()

	mouseCallback := func(w *glfw.Window,
		button glfw.MouseButton,
		action glfw.Action,
		mods glfw.ModifierKey) {
		// Mouse Pressed
		// fmt.Println("Mouse Callback, ", button, " ", action)
		if button == glfw.MouseButton1 {
			if action == glfw.Press {
				input.click = true
				input.start_dragging = true
				input.is_dragging = true

			} else if action == glfw.Release {
				input.is_dragging = false
				input.end_dragging = true
			}
		}
	}

	cursorCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		// Cursor Moved
		// fmt.Printf("Moved: %v, %v\n", xpos, ypos)
		input.current_mouse_x = xpos
		input.current_mouse_y = ypos
	}

	window.SetMouseButtonCallback(mouseCallback)
	window.SetCursorPosCallback(cursorCallback)

	for !window.ShouldClose() {
		startTime = glfw.GetTime()
		frameTime = startTime - lastStartTime
		lastStartTime = startTime

		// Do Opengl stuff
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT |
			gl.DEPTH_BUFFER_BIT |
			gl.STENCIL_BUFFER_BIT)

		C.nvgBeginFrame(vg, WindowWidth, WindowHeight, 2.0)

		UpdateAndRender(state, gui, frameTime)

		gameCodeTime = glfw.GetTime() - startTime

		// Display Stats
		stats := fmt.Sprintf("game: %.2fms, frame: %.2fms, fps: %.2f",
			f64round2(gameCodeTime*100.0),
			f64round2(frameTime*100.0),
			f64round2(1/frameTime))

		C.nvgSave(vg)
		C.nvgFontSize(vg, 14)
		C.nvgFillColor(vg, C.nvgRGBf(1, 1, 1))
		C.nvgText(vg, 5, WindowHeight-10, C.CString(stats), nil)
		C.nvgRestore(vg)

		// Show mouse position
		C.nvgSave(vg)
		C.nvgFillColor(vg, C.nvgRGBf(0.0, 1.0, 1.0))
		C.nvgBeginPath(vg)
		C.nvgCircle(vg, C.float(input.current_mouse_x), C.float(input.current_mouse_y), 2.0)
		C.nvgFill(vg)
		C.nvgRestore(vg)

		C.nvgEndFrame(vg)
		window.SwapBuffers()

		input.last_mouse_x = input.current_mouse_x
		input.last_mouse_y = input.current_mouse_y
		input.click = false
		input.start_dragging = false
		input.end_dragging = false
		glfw.PollEvents()
	}
}
