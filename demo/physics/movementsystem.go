package physics

import (
	"time"
	"fmt"

	"../../scene"
)

type MovementSystem struct {
	do chan func(*MovementSystem)
	running bool
	scene *scene.S
}

func New(s *scene.S) *MovementSystem {
	sys := &MovementSystem{do:make(chan func(*MovementSystem))}
	sys.Start(s)
	sys.Syn()
	return sys
}

func (sys *MovementSystem)Start(s *scene.S) {
	// bool check not thread safe
	if !sys.running {
		sys.scene = s

		// begin handling of events
		go func() {
			// no need to restrict to single thread (yet?)
			//runtime.LockOSThread()

			if err := sys.setup(); err != nil {
				panic(err.Error())
			}

			sys.running = true

			// start command handler
			for sys.running {
				visit := <-sys.do
				visit(sys)
			}

			sys.teardown()
		}()

		sys.Syn()

		// launch update timer
		ticker := time.Tick(16 * time.Millisecond)
		go func() {
			last := time.Now()
			for sys.running {
				now := <-ticker
				sys.tick(now.Sub(last))
				last = now
			}
		}()
	}
}

func (sys *MovementSystem)Stop() {
	if sys.running {
		sys.do <- func(sys *MovementSystem) {
			sys.running = false
		}
	}
}

func (sys *MovementSystem)tick(timestep time.Duration) {
	ts := timestep
	sys.do <- func(sys *MovementSystem) {
		sys.Update(ts)
	}
}

func (sys *MovementSystem)Syn() {
	ack := make(chan bool)
	sys.do <- func(sys *MovementSystem) {
		ack <- true
	}
	<-ack
	close(ack)
}

func (sys *MovementSystem)setup() error {
	// no setup needed right now
	println("setup movement system complete")
	return nil
}

func (sys *MovementSystem)teardown() {
	// no teardown needed right now
	println("teardown movement system complete")
}

func (sys *MovementSystem) Update(timestep time.Duration) (bool, error) {

	fmt.Printf("movement: %v\n", timestep)

	if actor := sys.scene.Actors["head"]; actor != nil {
		phy := actor.Get(PidPos).(*Pos)
		//phy.X += 0.01
		//phy.Y += 0.02
		phy.Z -= 0.1
	} else {
		fmt.Println("actor \"a\" not found")
	}


	return true,nil
}
