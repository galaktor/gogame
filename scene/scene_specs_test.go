package scene

import (
	"github.com/orfjackal/gospec"
	"testing"
)

func TestScene(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(CtorSpec)
	r.AddSpec(AddSpec)
	r.AddSpec(RemovePropertySpec)
	r.AddSpec(RemoveTypeSpec)
	r.AddSpec(FindSpec)

	gospec.MainGoTest(r, t)
}
