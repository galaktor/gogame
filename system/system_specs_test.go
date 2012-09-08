package system

import (
	"github.com/orfjackal/gospec"
	"testing"
)

func TestSystem(t *testing.T) {
	r := gospec.NewRunner()
	
	r.AddSpec(StartSpec)
	
	gospec.MainGoTest(r, t)
}