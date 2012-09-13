package graphics

import (
	"fmt"
	"time"
	"errors"

	"./gogre3d"
	"../physics"
	"../../scene"
)

type RenderSystem struct {
	s *scene.S
	// TODO: prop ids used

	r *gogre3d.Root
	w *gogre3d.RenderWindow
	Stop chan bool
}

func renderSystem(s *scene.S) *RenderSystem {
	return &RenderSystem{s:s,Stop:make(chan bool)}
}

func (r *RenderSystem)Init() error {
	root := gogre3d.NewRoot("./graphics/plugins.cfg", "./graphics/ogre.cfg", "ogre.log")
	r.r = &root

	if cancel := !r.r.ShowConfigDialog(); cancel {
		r.r.Delete()
		return errors.New("ogre dialog cancelled")
	}

	window := r.r.Initialise(true, "MyWindow")
	r.w = &window

	sm := r.r.CreateSceneManager()
	cam := sm.CreateCamera("MyCamera")
	cam.SetPosition(0, 0, 80)
	cam.LookAt(0, 0, -300)
	cam.SetNearClipDistance(5)

	vp := window.AddViewport(cam)
	vp.SetBackgroundColour(0, 0, 0, 0)
	
	w, h := vp.GetActualWidth(), vp.GetActualHeight()
	cam.SetAspectRatio(w / h)
	
	tm := gogre3d.GetTextureManager()
	tm.SetDefaultNumMipmaps(5)

	rm := gogre3d.GetResourceGroupManager()
	rm.AddResourceLocation("./graphics/gogre3d/Media/fonts", "FileSystem")
	rm.AddResourceLocation("./graphics/gogre3d/Media/models", "FileSystem")
	rm.AddResourceLocation("./graphics/gogre3d/Media/materials/scripts", "FileSystem")
	rm.AddResourceLocation("./graphics/gogre3d/Media/materials/programs", "FileSystem")
	rm.AddResourceLocation("./graphics/gogre3d/Media/materials/textures", "FileSystem")
	rm.InitialiseAllResourceGroups()

	head := sm.CreateEntity("Head", "ogrehead.mesh")
	rnode := sm.GetRootSceneNode()
	headnode := rnode.CreateChildSceneNode()
	headnode.AttachObject(head)

	sm.SetAmbientLight(0.5, 0.5, 0.5, 0)
	
	light := sm.CreateLight("MyLight")
	light.SetPosition(20, 80, 50)

	return nil
}

func (r *RenderSystem)Frequency() time.Duration {
	return 40 * time.Millisecond
}

func (sys *RenderSystem) Update(timestep time.Duration) (stop bool, err error) {
	fmt.Printf("render: %v\n", timestep)

	for _, actor := range sys.s.Find(physics.Pid, Pid) {
		p := actor.Get(physics.Pid).(*physics.Pos)
		g := actor.Get(Pid).(*Mesh)
		fmt.Printf("before push: %+v %+v\n", p, g)
		p.PushSync(g)
		fmt.Printf("after push: %+v %+v\n", p, g)
	}
	
	if err = sys.RenderOne(timestep); err != nil {
		stop = true
	}

	if !stop {
		stop = sys.checkStop()
	}
	return
}

func (sys *RenderSystem)checkStop() bool {
	select {
	case <-sys.Stop:
		println("received stop signal")
		return true
	default:
		return false
	}

	return false
}

func (sys *RenderSystem)RenderOne(timestep time.Duration) error {
	gogre3d.MessagePump()

	if sys.w.IsClosed() {
		return errors.New("window closed")
	}

	if !sys.r.RenderOneFrame() {
		return errors.New("error during render")
	}

	return nil
}

func (sys *RenderSystem)Exit() {
	sys.r.Delete()

	sys.Stop <- true
}