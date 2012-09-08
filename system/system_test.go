package system

import (
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
)

func StartSpec(c gospec.Context) {
	c.Specify("I NEED TESTS!", func() {
		c.Expect(true, Equals, false)
	})
}