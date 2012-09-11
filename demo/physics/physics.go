package physics

import (
	"time"
	"fmt"

	. "../../scene"
	"../../system"
)

var Pid = PropertyType(1)

type Physical struct {
	Do      chan func(*Physical)
	X, Y, Z float32
}

func Prop() *Physical {
	p := &Physical{}
	p.Do = make(chan func(*Physical))
	p.Start()
	return p
}

func (p *Physical) Type() PropertyType {
	return Pid
}

func (p *Physical) Start() {
	go func() {
		for {
			visit := <-p.Do
			visit(p)
		}
	}()
}

type CanPullPhysical interface {
	Pull(p *Physical)
}
type CanPushPhysical interface {
	Push(p CanPullPhysical)
}

func (p *Physical) Push(to CanPullPhysical) {
	to.Pull(p)
}

func (p *Physical) PushSync(to CanPullPhysical) {
	p.Do <- func(p *Physical) { p.Push(to) }
}

type MovementSystem struct {
	s *Scene
	// TODO: prop ids used
}

func movementSystem(s *Scene) *MovementSystem {
	return &MovementSystem{s}
}

func (m *MovementSystem)Frequency() time.Duration {
	return 16 * time.Millisecond
}

func (p *MovementSystem) Update(timestep time.Duration) {
	fmt.Printf("movement: %v\n", timestep)
	actor := p.s.Actors["a"]
	phy := actor.Get(Pid).(*Physical)
	phy.X += 1
	phy.Y += 2
	phy.Z += 3
}

func Start(s *Scene) *MovementSystem {
	sys := movementSystem(s)
	go system.Start(sys)
	return sys
}
