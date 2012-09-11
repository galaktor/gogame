package system

import (
	. "../scene"
	"fmt"
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
	"runtime"
	"time"
)

var s = NewScene()

func StartSpec(c gospec.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	a := s.Add(ActorId("a"))
	p := NewPhysical()
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