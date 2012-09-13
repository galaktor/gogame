package scene

import "testing"
import "fmt"

func BenchmarkFind(b *testing.B) {
	s := New()

	for i := 0; i < b.N; i++ {
		n := fmt.Sprintf("%v", i)

		a := s.Add(n)
		a.Add(NewProperty(PType(i)))

		
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Find(PType(b.N))
	}
}