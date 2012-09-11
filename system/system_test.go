package system

import (
	. "../scene"
	"fmt"
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
	"runtime"
	"time"
)

var propt = map[string]PropertyType{
	"phy": PropertyType(1),
	"gra": PropertyType(2),
}

type Graphical struct {
	Do      chan func(*Graphical)
	X, Y, Z float32
}

func (g *Graphical) Type() PropertyType {
	return propt["gra"]
}

func NewGraphical() *Graphical {
	g := &Graphical{}
	g.Do = make(chan func(*Graphical))
	g.Start()
	return g
}
func (g *Graphical) Start() {
	go func() {
		for {
			visit := <-g.Do
			visit(g)
		}
	}()
}

func (g *Graphical) Pull(p *Physical) {
	// capture variables
	x, y, z := p.X, p.Y, p.Z
	g.Do <- func(g *Graphical) {
		g.X += x
		g.Y += y
		g.Z += z
	}
}

type Physical struct {
	Do      chan func(*Physical)
	X, Y, Z float32
}

func NewPhysical() *Physical {
	p := &Physical{}
	p.Do = make(chan func(*Physical))
	p.Start()
	return p
}

func (p *Physical)Type() PropertyType {
	return propt["phy"]
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

type PhysicsSystem struct {
}

func (p *PhysicsSystem) Update(timestep time.Duration) {
	fmt.Printf("updating: %v\n", timestep)
	for _, actor := range s.Find(propt["phy"], propt["gra"]) {
		p := actor.Get(propt["phy"]).(*Physical)
		g := actor.Get(propt["gra"]).(*Graphical)
		fmt.Printf("got: %+v %+v\n", p, g)
		p.Push(g)
	}

}

var s = NewScene()

func StartSpec(c gospec.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	a := s.Add(ActorId("a"))
	p := NewPhysical()
	p.X, p.Y, p.Z = 1, 2, 3
	p.Start()
	a.Add(p)

	g := NewGraphical()
	g.Start()
	a.Add(g)
	
	sys := &PhysicsSystem{}
	
	Start(sys, 16 * time.Millisecond)

	c.Specify("I NEED TESTS!", func() {
		c.Expect(true, Equals, false)
	})
}