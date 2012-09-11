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
	graphics.Start(s)

	a := s.Add("a")
	a.Add(physics.Prop())
	a.Add(graphics.Prop())
	
	time.Sleep(5 * time.Second)
}