package physics

import (
	"../../scene"
)

var Pid = scene.PType(1)

type Pos struct {
	Do      chan func(*Pos)
	X, Y, Z float32
}

func NewPos() *Pos {
	p := &Pos{}
	p.Do = make(chan func(*Pos))
	p.start()
	return p
}

func (p *Pos) Type() scene.PType {
	return Pid
}

func (p *Pos) start() {
	go func() {
		for {
			visit := <-p.Do
			visit(p)
		}
	}()
}

type CanPullPosition interface {
	Pull(p *Pos)
}
type CanPushPosition interface {
	Push(p CanPullPosition)
}

func (p *Pos) Push(to CanPullPosition) {
	to.Pull(p)
}

func (p *Pos) PushSync(to CanPullPosition) {
	p.Do <- func(p *Pos) { p.Push(to) }
}