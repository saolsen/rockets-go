package main

import (
	// "errors"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

// #include "nanovg.h"
import "C"

func v2(x, y float32) mgl32.Vec2 {
	return mgl32.Vec2{x, y}
}

type BoundingBox struct {
	top_left     mgl32.Vec2
	bottom_right mgl32.Vec2
}

type NodeBounds struct {
	node_id     int
	position    BoundingBox
	input_index int
}

type NodeData struct {
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
	input    int
	NodeData
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
	NodeData
}

type Gate int

const (
	AND Gate = iota
	OR
	NOT
)

type GateNode struct {
	gate   Gate
	input1 int
	input2 int
	NodeData
}

type GameNode interface {
	Data() *NodeData
	GetText() string
	GetBounds(gui GuiState) (BoundingBox, string)
	Eval(ns NodeStore, pos_x int, pos_y int, rotation int) bool
	// Move(x int, y int)
	// SetInput(id int)
	// Move(x int, y int)
}

func (node ThrusterNode) Data() *NodeData {
	return &node.NodeData
}

func (node GateNode) Data() *NodeData {
	return &node.NodeData
}

func (node PredicateNode) Data() *NodeData {
	return &node.NodeData
}

const X_PADDING = 12.0
const Y_PADDING = 12.0

func (node ThrusterNode) GetText() string {
	return ""
}

func (node GateNode) GetText() string {
	switch node.gate {
	default:
		panic("Add new predicate here too")
	case (AND):
		return "AND"
	case (OR):
		return "OR"
	case (NOT):
		return "NOT"
	}
}

func (node PredicateNode) GetText() string {
	var signal string
	var predicate string

	switch node.signal {
	case (POS_X):
		signal = "pos-x"
	case (POS_Y):
		signal = "pos-y"
	case (ROTATION):
		signal = "rotation"
	}

	switch node.predicate {
	case (LT):
		predicate = "<"
	case (GT):
		predicate = ">"
	case (LEQT):
		predicate = "<="
	case (GEQT):
		predicate = ">="
	case (EQ):
		predicate = "=="
	case (NEQ):
		predicate = "<>"
	}

	return fmt.Sprintf("%v %v %v", signal, predicate, node.value)
}

func CalcBounds(gui GuiState, pos mgl32.Vec2, str string) BoundingBox {
	C.nvgSave(gui.vg)
	C.nvgFontSize(gui.vg, 14)

	bounds := make([]C.float, 4)

	C.nvgTextBounds(gui.vg,
		C.float(pos[0]+X_PADDING),
		C.float(pos[1]+Y_PADDING),
		C.CString(str), nil, &bounds[0])

	result := BoundingBox{
		v2(float32(bounds[0])-X_PADDING, float32(bounds[1])),
		v2(float32(bounds[2])+X_PADDING, float32(bounds[3])+2*Y_PADDING)}

	C.nvgRestore(gui.vg)
	return result
}

func (node ThrusterNode) GetBounds(gui GuiState) (BoundingBox, string) {
	return BoundingBox{node.pos, node.pos.Add(v2(60, 70))}, ""
}

func (node GateNode) GetBounds(gui GuiState) (BoundingBox, string) {
	str := node.GetText()
	return CalcBounds(gui, node.pos, str), str
}

func (node PredicateNode) GetBounds(gui GuiState) (BoundingBox, string) {
	str := node.GetText()
	return CalcBounds(gui, node.pos, str), str
}

func (node ThrusterNode) Eval(ns NodeStore, pos_x int, pos_y int, rotation int) bool {
	parent, exists := nsGetNode(ns, node.input)
	if exists {
		return parent.Eval(ns, pos_x, pos_y, rotation)
	} else {
		return false
	}
}

func (node GateNode) Eval(ns NodeStore, pos_x int, pos_y int, rotation int) bool {
	in1, exists1 := nsGetNode(ns, node.input1)
	val1 := false
	if exists1 {
		val1 = in1.Eval(ns, pos_x, pos_y, rotation)
	}

	in2, exists2 := nsGetNode(ns, node.input2)
	val2 := false
	if exists2 {
		val2 = in2.Eval(ns, pos_x, pos_y, rotation)
	}

	switch node.gate {
	case (OR):
		return val1 || val2
	case (AND):
		return val1 && val2
	default: // NOT
		return !val1
	}
}

func (node PredicateNode) Eval(ns NodeStore, pos_x int, pos_y int, rotation int) bool {
	var signal int
	switch node.signal {
	case (POS_X):
		signal = pos_x
	case (POS_Y):
		signal = pos_y
	default: // rotation
		signal = rotation
	}

	switch node.predicate {
	case (LT):
		return signal < node.value
	case (GT):
		return signal < node.value
	case (LEQT):
		return signal < node.value
	case (GEQT):
		return signal < node.value
	case (EQ):
		return signal < node.value
	default: // NEQ
		return signal < node.value
	}
}

type NodeStore struct {
	nextId int
	nodes  map[int]GameNode
}

func newNodeStore() NodeStore {
	return NodeStore{nextId: 1, nodes: make(map[int]GameNode)}
}

func nsAddThruster(ns *NodeStore, thruster Thruster) {
	ns.nodes[ns.nextId] = ThrusterNode{thruster, 0, NodeData{ns.nextId, mgl32.Vec2{0, 0}}}
	ns.nextId += 1
}

func nsAddGate(ns *NodeStore, gate Gate) {
	ns.nodes[ns.nextId] = GateNode{gate, 0, 0, NodeData{ns.nextId, mgl32.Vec2{0, 0}}}
	ns.nextId += 1
}

func nsAddPredicate(ns *NodeStore, signal Signal, predicate Predicate, value int) {
	ns.nodes[ns.nextId] =
		PredicateNode{signal, predicate, value, NodeData{ns.nextId, mgl32.Vec2{0, 0}}}
	ns.nextId += 1
}

func nsGetNode(ns NodeStore, id int) (node GameNode, exists bool) {
	value, exists := ns.nodes[id]
	return value, exists
}

type GameState struct {
	nextId int
	nodes  NodeStore
}

func Setup() *GameState {
	state := &GameState{nextId: 0, nodes: newNodeStore()}
	nsAddThruster(&state.nodes, BP)

	return state
}

func boundsContains(tlx, tly, brx, bry, x, y float32) bool {
	return (x > tlx && x < brx && y > tly && y < bry)
}

type ButtonState int

const (
	NAH ButtonState = iota
	CLICK
	HOVER
)

func guiButton(gui GuiState, x, y, width, height float32) bool {
	button := NAH

	if boundsContains(x, y, x+width, y+height,
		float32(gui.input.current_mouse_x), float32(gui.input.current_mouse_y)) {
		if gui.input.click {
			button = CLICK
		} else {
			button = HOVER
		}
	}

	C.nvgSave(gui.vg)
	switch button {
	case (NAH):
		C.nvgFillColor(gui.vg, C.nvgRGBf(0.0, 0.0, 1.0))
	case (CLICK):
		C.nvgFillColor(gui.vg, C.nvgRGBf(0.0, 1.0, 0.0))
	case (HOVER):
		C.nvgFillColor(gui.vg, C.nvgRGBf(1.0, 0.0, 0.0))
	}

	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg, C.float(x), C.float(y), C.float(width), C.float(height))
	C.nvgFill(gui.vg)

	C.nvgRestore(gui.vg)

	return button == CLICK
}

func drawTextBox(gui GuiState, bounds BoundingBox, txt string) {
	C.nvgSave(gui.vg)

	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg,
		C.float(bounds.top_left[0]),
		C.float(bounds.top_left[1]),
		C.float(bounds.bottom_right[0]-bounds.top_left[0]),
		C.float(bounds.bottom_right[1]-bounds.top_left[1]))
	C.nvgFillColor(gui.vg, C.nvgRGBf(0.5, 0.5, 0.5))
	C.nvgFill(gui.vg)

	C.nvgFontSize(gui.vg, 14)
	C.nvgFillColor(gui.vg, C.nvgRGBf(1, 1, 1))
	C.nvgText(gui.vg,
		C.float(bounds.top_left[0]+X_PADDING),
		C.float(bounds.top_left[1]+2*Y_PADDING),
		C.CString(txt),
		nil)
	C.nvgRestore(gui.vg)
}

func UpdateAndRender(state *GameState, gui GuiState, dt float64) {
	// Buttons to create nodes
	if guiButton(gui, 10, 10, 50, 25) {
		nsAddPredicate(&state.nodes, POS_X, EQ, 0)
	}

	if guiButton(gui, 70, 10, 50, 25) {
		nsAddGate(&state.nodes, AND)
	}

	if guiButton(gui, 130, 10, 50, 25) {
		nsAddThruster(&state.nodes, BOOST)
	}

	// Nodes
	bodies := make([]NodeBounds, 0)
	// inputs := make([]NodeBounds, 0)
	// inputs = make([]NodeBounds)
	// outputs = make([]NodeBounds)
	// constants = make([]NodeBounds)

	for id, node := range state.nodes.nodes {
		// data := node.Data()
		bounds, txt := node.GetBounds(gui)

		// body
		body := NodeBounds{id, bounds, 0}
		bodies = append(bodies, body)

		switch node := node.(type) {
		default:
			panic("Not all types accounted for here.")
		case PredicateNode:
			drawTextBox(gui, bounds, txt)
		case GateNode:
			drawTextBox(gui, bounds, txt)
		case ThrusterNode:
			fmt.Println(node)
		}

	}

	// Render
	// Space background!
	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg, 660, 10, 600, 700)
	C.nvgFillColor(gui.vg, C.nvgRGBf(0.25, 0.25, 0.25))
	C.nvgFill(gui.vg)

	C.nvgSave(gui.vg)
	// draw the scene
	C.nvgRestore(gui.vg)

}
