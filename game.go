package main

import "github.com/go-gl/mathgl/mgl32"

// #include "nanovg.h"
import "C"

// GUI nodes
type Node struct {
	id  int
	pos mgl32.Vec2
}

type Thruster int

const (
	BP Thruster = iota
	BS
	SP
	SS
	BOOST
)

type ThrusterNode struct {
	thruster Thruster
	Node
}

type Predicate int

const (
	LT Predicate = iota
	GT
	LEQT
	GEQT
	EQ
	NEQ
)

type Signal int

const (
	POS_X Signal = iota
	POS_Y
	ROTATION
)

type PredicateNode struct {
	signal    Signal
	predicate Predicate
	value     int
	Node
}

type Gate int

const (
	AND Gate = iota
	OR
	NOT
)

type GateNode struct {
	gate Gate
	Node
}

type NodeStore struct {
	nextId int
	nodes  map[int]Node
}

func newNodeStore() NodeStore {
	return NodeStore{nextId: 1, nodes: make(map[int]Node)}
}

type GameState struct {
	nextId int
	nodes  NodeStore
}

func Setup() *GameState {
	state := &GameState{nextId: 0, nodes: newNodeStore()}

	return state
}

func UpdateAndRender(vg *C.struct_NVGcontext, state *GameState, dt float64) {
	// state.pos[0] += 1 * float32(dt)

	// C.nvgBeginPath(vg)
	// C.nvgRect(vg, C.float(state.pos[0]), C.float(state.pos[1]), 10, 10)
	// C.nvgFillColor(vg, C.nvgRGBf(255, 255, 255))
	// C.nvgFill(vg)

	// mgl32.DegToRad(45.0)
	// C.nvgSave(vg)
	// C.nvgFontSize(vg, 24)
	// C.nvgFillColor(vg, C.nvgRGBf(1, 1, 1))
	// C.nvgText(vg, 10, 200, C.CString("Hello World"), nil)
	// // C.nvgRestore(vg)
}
