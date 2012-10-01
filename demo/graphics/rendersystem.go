package graphics

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/galaktor/gogre3d"

	"../../scene"
	"../physics"
)

type RenderSystem struct {
	do      chan func(*RenderSystem)
	running bool
	scene   *scene.S

	// ogre internals
	r  gogre3d.Root
	w  gogre3d.RenderWindow
	sm gogre3d.SceneManager
}

func (sys *RenderSystem) Syn() {
	ack := make(chan bool)
	sys.do <- func(sys *RenderSystem) {
		ack <- true
	}
	<-ack
	close(ack)
}

func New(s *scene.S) *RenderSystem {
	sys := &RenderSystem{do:make(chan func(*RenderSystem))}
	sys.Start(s)
	sys.Syn()
	return sys
}



func (sys *RenderSystem)Start(s *scene.S) {
	// bool check not thread safe
	if !sys.running {
		sys.scene = s

		// begin handling of events
		go func() {
			// IMPORTANT: ogre needs to stay on one thread
			runtime.LockOSThread()

			// do ogre setup
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
		ticker := time.Tick(40 * time.Millisecond)
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



// TODO: make synchronous, wait until teardown complete!
func (sys *RenderSystem) Stop() {
	if sys.running {
		sys.do <- func(sys *RenderSystem) {
			sys.running = false
		}
	}
}


func (sys *RenderSystem)tick(timestep time.Duration) {
	// is this a copied value, or not?
	ts := timestep
	sys.do <- func(sys *RenderSystem) {
		sys.Update(ts)
	}
}



func (r *RenderSystem)setup() error {
	println("setting up render system")

	r.r = gogre3d.NewRoot("", "", "ogre.log")
	
	// setup OpenGL
	r.r.LoadPlugin("RenderSystem_GL")
	rs := r.r.GetRenderSystemByName("OpenGL Rendering Subsystem")
	rs.SetConfigOption("Full Screen", "No")
	rs.SetConfigOption("VSync", "No")
	rs.SetConfigOption("Video Mode", "800 x 600 @ 32-bit")
	r.r.SetRenderSystem(rs)

	r.w = r.r.Initialise(true, "gogame demo")
	r.sm = r.r.CreateSceneManager("DefaultSceneManager", "The SceneManager")

	gogre3d.SetDefaultNumMipmaps(5)

	rgm := gogre3d.GetResourceGroupManager()
	rgm.AddResourceLocation("./graphics/media/models", "FileSystem", "Default")
	rgm.AddResourceLocation("./graphics/media/materials/scripts", "FileSystem", "Default")
	rgm.AddResourceLocation("./graphics/media/materials/textures", "FileSystem", "Default")
	rgm.InitialiseAllResourceGroups()

	println("setup render system complete")

	return nil
}

func (sys *RenderSystem)teardown() {
	println("tearing down render system")

	sys.r.Destroy()

	println("teardown render system complete")
}

func (sys *RenderSystem) Update(timestep time.Duration) (stop bool, err error) {
	fmt.Printf("render: %v\n", timestep)

	for _, actor := range sys.scene.Find(physics.PidPos, PidMesh) {
		p := actor.Get(physics.PidPos).(*physics.Pos)
		m := actor.Get(PidMesh).(*Mesh)
		m.ogre_update(p)
		fmt.Printf("mesh after push: %+v %+v\n", p, m)
	}

	for _, actor := range sys.scene.Find(physics.PidPos, PidCam) {
		p := actor.Get(physics.PidPos).(*physics.Pos)
		c := actor.Get(PidCam).(*Camera)
		c.ogre_update(p)
		fmt.Printf("cam after push: %+v %+v\n", p, c)
	}
	
	for _, actor := range sys.scene.Find(physics.PidPos, PidLight) {
		p := actor.Get(physics.PidPos).(*physics.Pos)
		l := actor.Get(PidLight).(*Light)
		l.ogre_update(p)
		fmt.Printf("light after push: %+v %+v\n", p, l)
	}

	if err = sys.renderOne(timestep); err != nil {
		return true, errors.New("error during ogre.RenderOne()")
	}

	return false, nil
}

func (sys *RenderSystem) renderOne(timestep time.Duration) error {
	gogre3d.MessagePump()

	if sys.w.IsClosed() {
		return errors.New("window closed")
	}

	if !sys.r.RenderOneFrameEx(float32(timestep.Nanoseconds())/1000000) {
		return errors.New("error during render")
	}

	return nil
}

func (sys *RenderSystem) Ambient(r, g, b float32) {
	sys.do <- func(sys *RenderSystem) {
		sys.sm.SetAmbientLight(r, g, b)
	}
}
