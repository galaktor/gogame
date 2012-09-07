package scene

import (
	"github.com/orfjackal/gospec"
	"testing"
)

func TestAllSpecs(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(SceneCtorSpec)
	r.AddSpec(SceneAddSpec)
	r.AddSpec(SceneRemoveByTypeSpec)
	r.AddSpec(SceneFindSpec)
	gospec.MainGoTest(r, t)
}
