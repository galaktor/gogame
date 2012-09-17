package graphics

import (
	"github.com/galaktor/gogre3d"

	"../physics"
	"../../scene"
)

var PidCam = scene.PType(3)

// rename to "View" or something
type Camera struct {
	do chan func(*Camera)

	// buffer changes
	lookx, looky, lookz float32
	nearclip float32

	// ogre internals; owned by render system
	// only to be used on ogre thread
	c gogre3d.Camera
	vp gogre3d.Viewport
}

func (c *Camera) Type() scene.PType {
	return PidCam
}

func (sys *RenderSystem) Cam(name string) *Camera {
	c := &Camera{}
	c.do = make(chan func(*Camera))

	c.start()

	println("pushing cam creation")
	sys.do <- func(sys *RenderSystem) {

		// create ogre stuff in rendersystem thread
		ogrec := sys.sm.CreateCamera(name)
		ogrevp := sys.w.AddViewport(ogrec)
		ogrevp.SetBackgroundColour(0, 0, 0, 0)
		w, h := ogrevp.GetActualWidth(), ogrevp.GetActualHeight()
		ogrec.SetAspectRatio(w / h)

		c.do <- func(m *Camera) {
			println("cam assigning ogre vars")
			c.c = ogrec
			c.vp = ogrevp

		}
	}
	
	sys.Syn()
	c.Syn()

	return c
}

func (c *Camera)Syn() {
	ack := make(chan bool)
	c.do <- func(c *Camera) {
		ack <- true
	}
	<-ack
	close(ack)
}

func (c *Camera)start() {
	go func() {
		for {
			visit := <-c.do
			visit(c)
		}
	}()
}


// flushes cached values into the ogre objects
// ONLY to be called from the ogre thread
func (c *Camera)ogre_update(p *physics.Pos) {
	c.c.SetPosition(p.X, p.Y, p.Z)
	c.c.LookAt(c.lookx, c.looky, c.lookz)
	c.c.SetNearClipDistance(c.nearclip)
}

// to be called from other goroutines
func (c *Camera) LookAt(x, y, z float32) {
	c.do <- func(c *Camera) {
		c.lookx = x
		c.looky = y
		c.lookz = z
	}
	
}

// to be called from other goroutines
func (c *Camera) NearClip(clip float32) {
	c.do <- func(c *Camera) {
		c.nearclip = clip
	}
}