/*
 Contains types and functions for systems.
*/
package system

import (
	"time"
	"runtime"
	"errors"
)

type System interface {
	DoOne()
	Syn()
	Running(bool)
	IsRunning()
	Frequency() time.Duration
	Setup() error
	Tick(time.Duration)
	TearDown()
}

func Start(sys System, lockOsThread bool) (err error) {
	// TODO: IsRunning check not atomic
	// checked here, set further down
	if sys.IsRunning() {
		return errors.New("system already running")

	} else {
		go func() {
			if lockOsThread {
				runtime.LockOSThread()
			}
			if e := sys.Setup(); e != nil {
				err = e
				return
			}

			sys.Running(true)
			for sys.IsRunning() {
				sys.DoOne()
			}

			sys.TearDown()
		}()

		sys.Syn()

		ticker := time.Tick(sys.Frequency())
		go func() {
			last := time.Now()
			for sys.IsRunning() {
				now := <-ticker
				sys.Tick(now.Sub(last))
				last = now
			}
		}()	
	}

	return nil
}
