/*
 Provides basic Physical properties for actors in the scene and a
 dummy MovementSystem that increments positional values.
*/
package physics

import (
	"../../scene"
	"../../system"
)

func Start(s *scene.S) *MovementSystem {
	sys := movementSystem(s)
	go system.Start(sys, false)
	return sys
}
