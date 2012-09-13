/*
 Contains types and functions for systems.
*/
package system

import (
	"time"
	"runtime"
	"errors"
)

// Every system must provide it's tick frequency as well
// as an Update method to be called each tick. Init
// and Exit allow for setup and teardown of the system
type System interface {
	Frequency() time.Duration
	Init() error
	Update(timestep time.Duration) (stop bool, err error)
	Exit()
}

// Starts a goroutine at the system's frequency which
// periodically calls it's update method, providing
// the "timestep", which is the time since the last
// update has completed.
// TODO: make timestep reflect if a system is falling behind
// it's frequency
func Start(s System, lockthread bool) (err error) {
	if lockthread {
		runtime.LockOSThread()
	}

	if e := s.Init(); e != nil {
		err = errors.New("error initialising system: " + e.Error())
		return
	}

	ticker := time.Tick(s.Frequency())
	last := time.Now()
	for now := range ticker {
		if stop,e := s.Update(now.Sub(last)); stop || e != nil{
			err = e
			
			s.Exit()
			break
		}
		
		last = now
	}

	return
}