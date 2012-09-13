package system

import (
	"errors"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
)

func StartSpec(c gospec.Context) {
	sys := NewSystem()

	c.Specify("When initialising a system", func() {
		c.Specify("and Init() returns nil", func() {
			sys.initReturn = nil
			err := Start(sys, false)

			c.Specify("Start() does not return an error", func() {
				c.Expect(err, IsNil)
			})

			c.Specify("Update() is called until it exits", func() {
				c.Expect(sys.updateCount, Equals, sys.maxUpdates)
			})

			c.Specify("Exit() is called once", func() {
				c.Expect(sys.exitCount, Equals, 1)
			})
		})

		c.Specify("and Init() returns an error", func() {
			sys.initReturn = errors.New("a fake error")
			err := Start(sys, false)

			c.Specify("Start() returns that error", func() {
				c.Expect(err, Satisfies, strings.Contains(err.Error(), "a fake error"))
			})

			c.Specify("Update() is never called", func() {
				c.Expect(sys.updateCount, Equals, 0)
			})

			c.Specify("Exit() is never called", func() {
				c.Expect(sys.exitCount, Equals, 0)
			})
		})
	})

	c.Specify("When updating a system", func() {
		c.Specify("timestep is within 1ms of the defined Frequency()", func() {
			sys.maxUpdates = 100
			Start(sys, false)

			c.Expect(sys.avgTimestepMs(), IsWithin(1.0), float64(sys.Frequency().Nanoseconds()/1000000))
		})

		c.Specify("and Update() returns true", func() {
			sys.maxUpdates = 2
			sys.updateReturnStop = true
			Start(sys, false)

			c.Specify("Update() is not called anymore after exit", func() {
				c.Expect(sys.updateCount, Equals, 1)
			})

			c.Specify("Exit() is called once", func() {
				c.Expect(sys.exitCount, Equals, 1)
			})
		})

		c.Specify("and Update() returns an error", func() {
			sys.maxUpdates = 2
			sys.updateReturnErr = errors.New("a fake error")
			err := Start(sys, false)

			c.Specify("Update() is not called anymore after exit", func() {
				c.Expect(sys.updateCount, Equals, 1)
			})

			c.Specify("Exit() is called once", func() {
				c.Expect(sys.exitCount, Equals, 1)
			})

			c.Specify("Start() returns that error", func() {
				c.Expect(err, Satisfies, strings.Contains(err.Error(), "a fake error"))
			})
		})
	})

	c.Specify("When \"lockthread\" param", func() {
		// we need lots of threads for this test to word
		prevprocs := runtime.GOMAXPROCS(runtime.NumCPU()*10)
		// start other yielding routines
		go func() { for { runtime.Gosched() }}()
		go func() { for { runtime.Gosched() }}()
		go func() { for { runtime.Gosched() }}()

		sys.maxUpdates = 1000
	
		c.Specify("is TRUE then Init(), Update() and Exit() only ever ran on one and the same thread", func() {
			Start(sys, true)

			c.Assume(len(sys.initTids), Equals, 1)
			c.Assume(len(sys.updateTids), Equals, 1)			
			c.Assume(len(sys.exitTids), Equals, 1)

			c.Specify("Init() thread equal to Update() thread", func() {
				c.Expect(sys.initTids[0], Equals, sys.updateTids[0])
			})

			c.Specify("Update() thread equal to Exit() thread", func() {
				c.Expect(sys.updateTids[0], Equals, sys.exitTids[0])
			})

		})

		c.Specify("is FALSE then Init(), Update() and Exit() only ever ran on one and the same thread", func() {

			Start(sys, false)

			c.Specify("either Init(), Update() or Exit() ran on more than one thread", func() {
				allTids := len(sys.initTids) + len(sys.updateTids) + len(sys.exitTids)
				c.Expect(allTids, Satisfies, allTids > 3)
			})
		})

		// reset to original thread number
		runtime.GOMAXPROCS(prevprocs)

	})
}



type FakeSystem struct {
	initTids, updateTids, exitTids []int
	initReturn                     error
	updateReturnStop               bool
	updateReturnErr                error
	maxUpdates, updateCount        int
	exitCount                      int
	totalTimeMs                    float64
}

func NewSystem() *FakeSystem {
	return &FakeSystem{maxUpdates: 1}
}

func (f *FakeSystem) Frequency() time.Duration {
	return 1 * time.Millisecond
}

func (f *FakeSystem) Init() error {
	f.initTids = AppendTid(f.initTids, syscall.Gettid())

	return f.initReturn
}

func (f *FakeSystem) avgTimestepMs() float64 {
	return f.totalTimeMs / float64(f.updateCount)
}

func AppendTid(list []int, tid int) []int {
	found := false
	for _, id := range list {
		if id == tid {
			found = true
			break
		}
	}

	if !found {
		return append(list, tid)
	}

	return list
}

func (f *FakeSystem)Update(timestep time.Duration) (stop bool, err error) {
	f.updateTids = AppendTid(f.updateTids, syscall.Gettid())

	f.updateCount++
	if f.updateCount >= f.maxUpdates {
		stop = true
	}

	f.totalTimeMs += float64(timestep.Nanoseconds())/1000000
	
	runtime.Gosched()


	stop = stop || f.updateReturnStop
	err = f.updateReturnErr
	return
}

func (f *FakeSystem)Exit() {
	f.exitTids = AppendTid(f.exitTids, syscall.Gettid())

	f.exitCount++
}