package main

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// #include "nanovg.h"
import "C"

func v2(x, y float32) mgl32.Vec2 {
	return mgl32.Vec2{x, y}
}

type Thrusters struct {
	bp    bool
	bs    bool
	sp    bool
	ss    bool
	boost bool
}

type Ship struct {
	position  mgl32.Vec2
	rotation  int
	thrusters Thrusters
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
	default:
		panic("Add new signal here too")
	case (POS_X):
		signal = "pos-x"
	case (POS_Y):
		signal = "pos-y"
	case (ROTATION):
		signal = "rotation"
	}

	switch node.predicate {
	default:
		panic("Add new predicate here too")
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

func (node ThrusterNode) Eval(ns NodeStore, pos_x, pos_y, rotation int) bool {
	parent, exists := nsGetNode(ns, node.input)
	if exists {
		return parent.Eval(ns, pos_x, pos_y, rotation)
	} else {
		return false
	}
}

func (node GateNode) Eval(ns NodeStore, pos_x, pos_y, rotation int) bool {
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

func (node PredicateNode) Eval(ns NodeStore, pos_x, pos_y, rotation int) bool {
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

type LevelStatus int

const (
	RUNNING LevelStatus = iota
	PAUSED
	WON
	DIED
)

type GameState struct {
	// nodes
	nodes NodeStore

	// scene
	ship         Ship
	currentLevel int
	goal         mgl32.Vec2
	status       LevelStatus
}

func Setup() *GameState {
	state := &GameState{
		nodes: newNodeStore(),
		ship: Ship{
			position: v2(300, 99),
			rotation: 0,
		},
		currentLevel: 1,
		goal:         v2(300, 600),
		status:       RUNNING,
	}

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

func drawShip(gui GuiState, thrusters Thrusters, grayscale bool) {
	C.nvgSave(gui.vg)

	if grayscale {
		C.nvgFillColor(gui.vg, C.nvgRGBf(1, 1, 1))
	} else {
		C.nvgFillColor(gui.vg, C.nvgRGBf(1.0, 0.0, 0.0))
	}

	// Ship body
	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg, -10.0, -25.0, 20.0, 40.0)
	C.nvgFill(gui.vg)

	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg, -20.0, -5.0, 15.0, 30.0)
	C.nvgFill(gui.vg)

	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg, 5.0, -5.0, 15.0, 30.0)
	C.nvgFill(gui.vg)

	// Ship thrusters
	if grayscale {
		C.nvgFillColor(gui.vg, C.nvgRGBf(0.75, 0.75, 0.75))
	} else {
		C.nvgFillColor(gui.vg, C.nvgRGBf(1.0, 1.0, 0.0))
	}

	if thrusters.bp {
		C.nvgBeginPath(gui.vg)
		C.nvgRect(gui.vg, -20.0, -25.0, 10, 10)
		C.nvgFill(gui.vg)
	}

	if thrusters.bs {
		C.nvgBeginPath(gui.vg)
		C.nvgRect(gui.vg, 10.0, -25.0, 10, 10)
		C.nvgFill(gui.vg)
	}

	if thrusters.sp {
		C.nvgBeginPath(gui.vg)
		C.nvgRect(gui.vg, -30.0, 15.0, 10, 10)
		C.nvgFill(gui.vg)
	}

	if thrusters.ss {
		C.nvgBeginPath(gui.vg)
		C.nvgRect(gui.vg, 20.0, 15.0, 10, 10)
		C.nvgFill(gui.vg)
	}

	if thrusters.boost {
		C.nvgBeginPath(gui.vg)
		C.nvgRect(gui.vg, -17.5, 25.0, 10, 10)
		C.nvgFill(gui.vg)
		C.nvgBeginPath(gui.vg)
		C.nvgRect(gui.vg, 7.5, 25.0, 10, 10)
		C.nvgFill(gui.vg)
	}

	C.nvgRestore(gui.vg)
}

func evalThrusters(nodes NodeStore, ship Ship) Thrusters {
	outThrusters := Thrusters{}

	for _, node := range nodes.nodes {
		// @TODO should I round instead of cast?
		value := node.Eval(nodes, int(ship.position[0]), int(ship.position[1]), ship.rotation)

		switch node := node.(type) {
		case (ThrusterNode):
			switch node.thruster {
			case (BP):
				outThrusters.bp = outThrusters.bp || value
			case (BS):
				outThrusters.bs = outThrusters.bs || value
			case (SP):
				outThrusters.sp = outThrusters.sp || value
			case (SS):
				outThrusters.ss = outThrusters.ss || value
			case (BOOST):
				outThrusters.boost = outThrusters.boost || value
			}
		}
	}

	return outThrusters
}

func moveShip(ship *Ship, dt float64) {
	force := v2(0, 0)
	rotation := 0

	if ship.thrusters.bp {
		force = force.Add(v2(1, 0))
		rotation--
	}

	if ship.thrusters.bs {
		force = force.Add(v2(-1, 0))
		rotation++
	}

	if ship.thrusters.sp {
		force = force.Add(v2(1, 0))
		rotation++
	}

	if ship.thrusters.ss {
		force = force.Add(v2(-1, 0))
		rotation--
	}

	if ship.thrusters.boost {
		force = force.Add(v2(0, 5))
	}

	// @TODO: figure out what to do about all these conversions.
	// Either do all math in 64 which is fine for this game because there is so little
	// or find a float32 math library.
	theta := float64(ship.rotation) * math.Pi / 180.0
	nx := float64(force[0])*math.Cos(theta) - float64(force[1])*math.Sin(theta)
	ny := float64(force[0])*math.Sin(theta) + float64(force[1])*math.Cos(theta)

	absForce := v2(float32(nx), float32(ny))
	speed := float32(50.0 * dt)

	ship.position = ship.position.Add(absForce.Mul(speed))

	newRotation := ship.rotation + rotation
	if newRotation < 0 {
		newRotation += 360
	}
	ship.rotation = newRotation % 360
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
	// inputs = make([]NodeBounds, 0)
	// outputs = make([]NodeBounds, 0)
	// constants = make([]NodeBounds, 0)

	for id, node := range state.nodes.nodes {
		bounds, txt := node.GetBounds(gui)

		// body
		body := NodeBounds{id, bounds, 0}
		bodies = append(bodies, body)

		switch node := node.(type) {
		default:
			panic("Not all node types accounted for here.")
		case PredicateNode:
			drawTextBox(gui, bounds, txt)
		case GateNode:
			drawTextBox(gui, bounds, txt)
		case ThrusterNode:
			thrusts := Thrusters{}

			switch node.thruster {
			case (BP):
				thrusts.bp = true
			case (BS):
				thrusts.bs = true
			case (SP):
				thrusts.sp = true
			case (SS):
				thrusts.ss = true
			case (BOOST):
				thrusts.boost = true
			}

			C.nvgSave(gui.vg)
			C.nvgBeginPath(gui.vg)
			C.nvgRect(gui.vg,
				C.float(bounds.top_left[0]),
				C.float(bounds.top_left[1]),
				60, 70)
			C.nvgFillColor(gui.vg, C.nvgRGBf(0.5, 0.5, 0.5))
			C.nvgFill(gui.vg)

			C.nvgTranslate(gui.vg,
				C.float(bounds.top_left[0]+30),
				C.float(bounds.top_left[1]+35))

			drawShip(gui, thrusts, true)

			C.nvgRestore(gui.vg)
		}
	}

	// Update
	if state.status == RUNNING {
		newThrusters := evalThrusters(state.nodes, state.ship)
		state.ship.thrusters = newThrusters
		moveShip(&state.ship, dt)
	}

	// Render

	// Space background!
	C.nvgBeginPath(gui.vg)
	C.nvgRect(gui.vg, 660, 10, 600, 700)
	C.nvgFillColor(gui.vg, C.nvgRGBf(0.25, 0.25, 0.25))
	C.nvgFill(gui.vg)

	// draw the scene
	C.nvgSave(gui.vg)

	// @HARDCODE
	C.nvgTranslate(gui.vg, 660, 10)
	C.nvgTranslate(gui.vg, 0, 700)

	// Must draw with negative x coordinates in this transform.
	// Uses 4th quadrant instead of 1st because that way text will render correctly.
	C.nvgSave(gui.vg)
	C.nvgTranslate(gui.vg,
		C.float(state.ship.position[0]),
		C.float(-state.ship.position[1]))
	C.nvgRotate(gui.vg, C.float(-(float64(state.ship.rotation) * math.Pi / 180.0)))
	drawShip(gui, state.ship.thrusters, false)
	C.nvgRestore(gui.vg)

	C.nvgRestore(gui.vg)

}
