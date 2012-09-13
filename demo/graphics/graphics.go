/*
 Provides dummy Graphical property for actors and a system
 that flushes positional data from Physical properties into
 the Graphical properties before (hypothetically) rendering
 a frame.
*/
package graphics

import (
	"../../scene"
	"../../system"
)

func Start(s *scene.S) *RenderSystem {
	sys := renderSystem(s)
	go func() {
		if err := system.Start(sys, true); err != nil {
			println("Start error: " + err.Error())
		}
	}() 
	return sys
}