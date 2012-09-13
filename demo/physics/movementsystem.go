package physics

import (
	"time"
	"fmt"

	"../../scene"
)

type MovementSystem struct {
	s *scene.S
	// TODO: prop ids used
}

func movementSystem(s *scene.S) *MovementSystem {
	return &MovementSystem{s}
}

func (m *MovementSystem)Frequency() time.Duration {
	return 16 * time.Millisecond
}

func (p *MovementSystem) Update(timestep time.Duration) (bool, error) {
	fmt.Printf("movement: %v\n", timestep)

	if actor := p.s.Actors["a"]; actor != nil {
		phy := actor.Get(Pid).(*Pos)
		phy.X += 1
		phy.Y += 2
		phy.Z += 3
	} else {
		fmt.Println("actor \"a\" not found")
	}

	return true,nil

}

func (m *MovementSystem)Init() error {
	return nil
}

func (m *MovementSystem)Exit() {
	// do nothing
}