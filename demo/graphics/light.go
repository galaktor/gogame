package graphics

import (
	"github.com/galaktor/gogre3d"

	"../physics"
	"../../scene"
)

var PidLight = scene.PType(4)

type Light struct {
	do chan func(*Light)

	// ogre internals; owned by render system
	// only to be used on ogre thread
	l gogre3d.Light
}

func (l *Light) Type() scene.PType {
	return PidLight
}

func (r *RenderSystem) Light(name string) *Light {
	l := &Light{}
	l.do = make(chan func(*Light))

	l.start()

	// force creation of light on render system thread
	// need nicer way to express this in API
	// compare to Push Pull model
	// basically the same
	println("pushing light creation")
	r.do <- func(r *RenderSystem) {
		l.do <- func(l *Light) {
			l.l = r.sm.CreateLight(name)
		}
		
	}
	

	return l
}

/*
func(l *Light)Do(cmd func(*Light)) {
	l.do <- cmd
}
*/

func (l *Light)start() {
	go func() {
		for {
			visit := <-l.do
			visit(l)
		}
	}()
}

func (l *Light)ogre_update(p *physics.Pos) {
	l.l.SetPosition(p.X, p.Y, p.Z)
}
