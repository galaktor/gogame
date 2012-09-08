package system

import "time"

type System interface {
	Update(timestep time.Duration)
}

func Start(s System, interval time.Duration) {
	ticker := time.Tick(interval)
	last := time.Now()
	for now := range ticker {
		s.Update(now.Sub(last))
		last = now
	}
}