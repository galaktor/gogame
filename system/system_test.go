package system

import (
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
	. "../scene"
	"fmt"
	"time"
	"runtime"
)

// should a system "have a" scene?
// how to keep orthogonal?
type MySystem struct {
	Scene
}

// TODO: need a nicer way to do this
func GetProperties(s *Scene, a Actor) (p0, p1, p2 *MyProperty) {
	for _,prop := range s.Properties[a] {
		switch prop.Type() {
		case 0:
			p0 = prop.(*MyProperty)
		case 1:
			p1 = prop.(*MyProperty)
		case 2:
			p2 = prop.(*MyProperty)
		}
	}

	return
}

func (sys *MySystem)Update(timestep time.Duration) {
	// TODO: optimize by performing Find() only when relevant props have been changed
	actors := sys.Scene.Find(PropertyType(Prop0),PropertyType(Prop1),PropertyType(Prop2))
	for _,a := range actors {
		p0,p1,p2 := GetProperties(&sys.Scene, a)
		
		fmt.Printf("%v %v props:\np0:%+v\np1:%+v\np2:%+v\n",timestep, a, p0, p1, p2) 

		p2.Pull(p1.pos)
	}
}

func (p *Vector3)Add(other Vector3) {
	p.x += other.x
	p.y += other.y
	p.z += other.z
}

func (to *MyProperty)Pull(pos Vector3) {
	to.In <- func(p *MyProperty) {
		p.pos.Add(pos)
	}
}

type Vector3 struct {
	x, y, z float32
}

type MyProperty struct {
	tid PropertyType
	pos Vector3
	In chan func(*MyProperty)
}

func (p *MyProperty)Type() PropertyType {
	return p.tid
}

func (p *MyProperty)Start() {
	go func() {
		for {
			f := <- p.In
			fmt.Printf("In @ %+v\n", p)
			f(p)
		}
	}()
}

func NewProperty(tid PropertyType) Property {
	p := &MyProperty{tid,Vector3{1, 2, 3},make(chan func(*MyProperty))}
	p.Start()
	return p
}


const (
	Prop0 = iota
	Prop1
	Prop2
)

func StartSpec(c gospec.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	s := NewScene()

	s.Add("a", NewProperty(Prop0))
	s.Add("a", NewProperty(Prop1))
	s.Add("a", NewProperty(Prop2))
	
	s.Add("b", NewProperty(Prop0))
	s.Add("b", NewProperty(Prop1))
	s.Add("b", NewProperty(Prop2))

	sys := &MySystem{*s}
	// in this example: min duration is ~100us
	Start(sys, 16 * time.Millisecond)

	c.Specify("I NEED TESTS!", func() {
		c.Expect(true, Equals, false)
	})
}