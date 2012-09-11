package main

import (
	"time"
	"runtime"
	"fmt"

	"./physics"
	"./graphics"
	"../scene"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// stutter!
	s := scene.NewScene()

	// get rid of ActorId alias
	a := s.Add(scene.ActorId("a"))

	// make them start automatically
	p := physics.Prop()
	g := graphics.Prop()

	a.Add(p)
	a.Add(g)

	// better: physics.System()
	// even better: physics.Start()
	physics.Start(s)
	graphics.Start(s)
	
	fmt.Println("sleeping")
	time.Sleep(5 * time.Second)
}