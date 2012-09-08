package system

import (
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
	"../scene")

type MySystem struct {

}

type MyProperty struct {

}

func (p *MyProperty)Type() scene.PropertyType {
	return scene.PropertyType(1)
}




func StartSpec(c gospec.Context) {
	s := scene.NewScene()

	s.Add("a", &MyProperty{})


	c.Specify("I NEED TESTS!", func() {
		c.Expect(true, Equals, false)
	})
}