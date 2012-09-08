package scene

import "testing"
import "fmt"

func BenchmarkFind(b *testing.B) {
	s := NewScene()

	for i := 0; i < b.N; i++ {
		n := fmt.Sprintf("%v", i)
		s.Add(Actor(n), NewProperty(PropertyType(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Find(PropertyType(b.N))
	}
}