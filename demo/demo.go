/*
 A simple demonstration of how to create a scene and launch systems.
*/
package main

import (
	"time"
	"runtime"

	"./physics"
	"./graphics"
	"../scene"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	s := scene.New()

	println("started")

	psys := physics.New(s)
	rsys := graphics.New(s)
	
	println("adding head")
	h := s.Add("head")
	h.Add(psys.Pos()) // 0, 0, 0
	println("create mesh")
	h.Add(rsys.Mesh("head", "ogrehead.mesh"))
	
	c := s.Add("camera")
	pp := psys.Pos()
	pp.Set(0, 0, 80)
	c.Add(pp)
	cp := rsys.Cam(c.Id)
	cp.LookAt(0, 0, -300)
	cp.NearClip(5)
	c.Add(cp)

	l := s.Add("light")
	pp = psys.Pos()
	pp.Set(20, 80, 50)
	l.Add(pp)
	l.Add(rsys.Light(l.Id))

	rsys.Ambient(0.5, 0.5, 0.5, 0)

	time.Sleep(3 * time.Second)

	// TODO: make synchronous
	rsys.Stop()
	time.Sleep(1 * time.Second)
}