/*
 Contains types and functions for systems.
*/
package system

import "time"

// Every system must provide it's tick frequency as well
// as an Update method which will be called every tick
type System interface {
	Frequency() time.Duration
	Update(timestep time.Duration)
}

// Starts a goroutine at the system's frequency which
// periodically calls it's update method, providing
// the "timestep", which is the time since the last
// update has completed.
// TODO: make timestep reflect if a system is falling behind
// it's frequency
func Start(s System) {
	ticker := time.Tick(s.Frequency())
	last := time.Now()
	for now := range ticker {
		s.Update(now.Sub(last))
		last = now
	}
}