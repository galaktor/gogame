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

	physics.Start(s)
	render := graphics.Start(s)

	a := s.Add("a")
	a.Add(physics.NewPos())
	a.Add(graphics.NewMesh())
	
	time.Sleep(10 * time.Second)

	render.Stop <- true
	println("waiting for render system to exit")
	<-render.Stop
}