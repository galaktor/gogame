package graphics

import (
	"github.com/galaktor/gogre3d"

	"../physics"
	"../../scene"
)

var PidMesh = scene.PType(2)

type Mesh struct {
	do chan func(*Mesh)

	// ogre internals; owned by render system
	// only to be used on ogre thread
	n gogre3d.SceneNode
	e gogre3d.Entity
}

func (m *Mesh) Type() scene.PType {
	return PidMesh
}

func (sys *RenderSystem) Mesh(name, src string) *Mesh {
	m := &Mesh{}
	m.do = make(chan func(*Mesh))

	m.start()

	println("pushing mesh creation")
	sys.do <- func(sys *RenderSystem) {

		// create ogre stuff in rendersystem thread
		ogree := sys.sm.CreateEntity(name, src)
		rnode := sys.sm.GetRootSceneNode()
		ogren := rnode.CreateChildSceneNode()
		// should be an action on the mesh property?
		ogren.AttachObject(ogree)

		m.do <- func(m *Mesh) {
			m.e = ogree
			m.n = ogren
		}
	}

	sys.Syn()
	m.Syn()

	return m
}

func (m *Mesh)start() {
	go func() {
		for {
			visit := <-m.do
			visit(m)
		}
	}()
}

func (m *Mesh)Syn() {
	ack := make(chan bool)
	m.do <- func(c *Mesh) {
		ack <- true
	}
	<-ack
	close(ack)
}

/*
func (m *Mesh)AttachTo(o *Mesh) {
	m.Detach()

	// TODO: consider; no child node, attach directly to o.n?
	m.n = o.n.CreateChildSceneNode()
	m.n.AttachObject(m.e)
}

func (m *Mesh)Detach() {
	// TODO: need nil check; gogre3d pointers necessary for nil check
	m.n.DetachObject(m.e)
	//// destroy is probably bad, may be other entities on that node
	//m.sm.Destroy(m.n)
}
*/

func (m *Mesh)ogre_update(p *physics.Pos) {
	m.n.SetPosition(p.X, p.Y, p.Z)
}