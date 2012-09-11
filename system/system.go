package system

import "time"

type System interface {
	Frequency() time.Duration
	Update(timestep time.Duration)
}

func Start(s System) {
	ticker := time.Tick(s.Frequency())
	last := time.Now()
	for now := range ticker {
		s.Update(now.Sub(last))
		last = now
	}
}