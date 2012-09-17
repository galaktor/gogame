package physics

import (
	"../../scene"
)

var PidPos = scene.PType(1)

type Pos struct {
	Do      chan func(*Pos)
	X, Y, Z float32
}

func (m *MovementSystem) Pos() *Pos {
	p := &Pos{}
	p.Do = make(chan func(*Pos))
	p.start()
	return p
}

func (p *Pos) Type() scene.PType {
	return PidPos
}

func (p *Pos) start() {
	go func() {
		for {
			visit := <-p.Do
			visit(p)
		}
	}()
}

func (p *Pos) Set(x, y, z float32) {
	// capture values
	x1, y1, z1 := x, y, z

	p.Do <- func(p *Pos) {
		p.X, p.Y, p.Z = x1, y1, z1
	}
}

type CanPullPosition interface {
	Pull(p *Pos)
}
type CanPushPosition interface {
	Push(p CanPullPosition)
}

func (p *Pos) Push(to CanPullPosition) {
	p.Do <- func(p *Pos) { to.Pull(p) }
}