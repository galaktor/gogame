package graphics

import (
	"../physics"
	"../../scene"
)

var Pid = scene.PType(2)

type Mesh struct {
	Do chan func(*Mesh)
	X, Y, Z float32
}

func (m *Mesh) Type() scene.PType {
	return Pid
}

func NewMesh() *Mesh {
	m := &Mesh{}
	m.Do = make(chan func(*Mesh))
	m.start()
	return m
}
func (m *Mesh) start() {
	go func() {
		for {
			visit := <-m.Do
			visit(m)
		}
	}()
}

func (m *Mesh) Pull(p *physics.Pos) {
	println("mesh.Pull")

	// capture variables
	x, y, z := p.X, p.Y, p.Z

	m.Do <- func(m *Mesh) {
		m.X = x
		m.Y = y
		m.Z = z
	}
}