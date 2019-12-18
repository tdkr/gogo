package model

import (
	"bytes"
	"fmt"
)

type Vector2 struct {
	X int32
	Y int32
}

func NewVec2(x, y int32) *Vector2 {
	return &Vector2{X: x, Y: y}
}

func (v *Vector2) HashCode() int32 {
	return v.X<<16 + v.Y<<8
}

func (v *Vector2) Set(x, y int32) {
	v.X = x
	v.Y = y
}

func (v *Vector2) SetX(x int32) {
	v.X = x
}

func (v *Vector2) SetY(y int32) {
	v.Y = y
}

func (v *Vector2) Equals(vec *Vector2) bool {
	return v.X == vec.X && v.Y == vec.Y
}

func (v *Vector2) String() string {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte('{')
	buffer.WriteString(fmt.Sprintf("x=%d,y=%d", v.X, v.Y))
	buffer.WriteByte('}')
	return buffer.String()
}

// NewVecStack returns a new stack.
func NewVecStack() *VecStack {
	return &VecStack{}
}

// VecStack is a basic LIFO stack that resizes as needed.
type VecStack struct {
	nodes []*Vector2
	count int
}

// Push adds a node to the stack.
func (s *VecStack) Push(n *Vector2) {
	s.nodes = append(s.nodes, n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *VecStack) Pop() *Vector2 {
	if s.count == 0 {
		return nil
	}
	s.count--
	s.nodes = s.nodes[:len(s.nodes)]
	return s.nodes[s.count]
}

func (s *VecStack) Find(n *Vector2) int {
	for i, v := range s.nodes {
		if v.Equals(n) {
			return i
		}
	}
	return -1
}

func (s *VecStack) Nodes() []*Vector2 {
	return s.nodes
}

func (s *VecStack) Size() int {
	return s.count
}

func (s *VecStack) String() string {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte('{')
	for i := 0; i < s.count; i++ {
		buffer.WriteString(s.nodes[i].String())
		buffer.WriteByte(',')
	}
	buffer.WriteByte('}')
	return buffer.String()
}
