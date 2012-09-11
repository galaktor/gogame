/*
 Provides dummy Graphical property for actors and a system
 that flushes positional data from Physical properties into
 the Graphical properties before (hypothetically) rendering
 a frame.
*/
package graphics

import (
	"fmt"
	"time"
	
	"../physics"
	. "../../scene"
	"../../system"
)

var Pid = PropertyType(2)

type Graphical struct {
	Do      chan func(*Graphical)
	X, Y, Z float32
}

func (g *Graphical) Type() PropertyType {
	return Pid
}

func Prop() *Graphical {
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

func (g *Graphical) Pull(p *physics.Physical) {
	fmt.Println("graphics.Pull")
	// capture variables
	x, y, z := p.X, p.Y, p.Z
	g.Do <- func(g *Graphical) {
		g.X = x
		g.Y = y
		g.Z = z
	}
}

type RenderSystem struct {
	s *Scene
	// TODO: prop ids used
}

func renderSystem(s *Scene) *RenderSystem {
	return &RenderSystem{s}
}

func (r *RenderSystem)Frequency() time.Duration {
	return 32 * time.Millisecond
}

func (sys *RenderSystem) Update(timestep time.Duration) {
	fmt.Printf("render: %v\n", timestep)
	// update mesh positions
	for _, actor := range sys.s.Find(physics.Pid, Pid) {
		p := actor.Get(physics.Pid).(*physics.Physical)
		g := actor.Get(Pid).(*Graphical)
		fmt.Printf("before push: %+v %+v\n", p, g)
		p.PushSync(g)
		fmt.Printf("after push: %+v %+v\n", p, g)
	}

	// render
	// TODO: render one frame
	fmt.Println("render one frame")
}

func Start(s *Scene) *RenderSystem {
	sys := renderSystem(s)
	go system.Start(sys)
	return sys
}