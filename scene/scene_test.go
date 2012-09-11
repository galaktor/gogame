package scene

import (
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
	"strings"
)

type SomeProperty struct {
	tid PropertyType
}

func NewProperty(t PropertyType) Property {
	return &SomeProperty{t}
}

func (p *SomeProperty) Type() PropertyType {
	return p.tid
}

func CtorSpec(c gospec.Context) {
	c.Specify("New scene", func() {
		scene := New()

		c.Specify("has zero actors", func() {
			c.Expect(len(scene.Actors), Equals, 0)
		})

		c.Specify("has zero cached by property", func() {
			c.Expect(len(scene.byProperty), Equals, 0)
		})
	})
}

func AddSpec(c gospec.Context) {
	scene := New()

	c.Specify("New actor always has reference to creating scene", func() {
		a := scene.Add("a")
		
		c.Expect(a.s, IsSame, scene)
	})

	c.Specify("Add one actor with one property", func() {
		p1 := NewProperty(1)
		a := scene.Add("a")
		a.Add(p1)

		c.Specify("scene contains just that one actor", func() {
			c.Expect(len(scene.Actors), Equals, 1)
			c.Expect(scene.Actors["a"], IsSame, a)

			c.Expect(len(scene.byProperty), Equals, 1)
			c.Expect(scene.byProperty[1], Contains, a)
		})

		c.Specify("actor contains that property", func() {
			c.Expect(len(a.properties), Equals, 1)
			c.Expect(a.properties[1], IsSame, p1)
		})
	})

	c.Specify("Add two properties of same type to one actor", func() {
		a := scene.Add("a")
		p1, p2 := NewProperty(42), NewProperty(42)
		a.Add(p1)
		e := a.Add(p2)

		c.Specify("returns error contains the duplicate property type id", func() {
			c.Expect(e, Not(IsNil))
			c.Expect(e, Satisfies, strings.Contains(e.Error(), "42"))
		})
	})

	c.Specify("Add two actors with different properties", func() {
		p1, p2 := NewProperty(1), NewProperty(2)
		a,b := scene.Add("a"),scene.Add("b")
		a.Add(p1)
		b.Add(p2)

		c.Specify("both actors contain those properties", func() {
			c.Expect(len(a.properties), Equals, 1)
			c.Expect(a.properties[1], IsSame, p1)
			c.Expect(len(b.properties), Equals, 1)
			c.Expect(b.properties[2], IsSame, p2)
		})

		c.Specify("scene contains both actors", func() {
			c.Expect(len(scene.Actors), Equals, 2)
			c.Expect(scene.Actors["a"], IsSame, a)
			c.Expect(scene.Actors["b"], IsSame, b)
			c.Expect(len(scene.byProperty), Equals, 2)
			c.Expect(scene.byProperty[1], Contains, a)
			c.Expect(scene.byProperty[2], Contains, b)
		})
	})

	c.Specify("Add two properties to same actor", func() {
		p1,p2 := NewProperty(1),NewProperty(2)
		a := scene.Add("a")
		a.Add(p1)
		a.Add(p2)

		c.Specify("actor contains both properties", func() {
			c.Expect(len(a.properties), Equals, 2)
			c.Expect(a.properties[1], IsSame, p1)
			c.Expect(a.properties[2], IsSame, p2)
		})

		c.Specify("scene contains that actor", func() {
			c.Expect(len(scene.Actors), Equals, 1)
			c.Expect(scene.Actors["a"], IsSame, a)
			c.Expect(len(scene.byProperty), Equals, 2)
			c.Expect(scene.byProperty[1], Contains, a)
			c.Expect(scene.byProperty[2], Contains, a)
		})
	})

	c.Specify("Add same property to two actors", func() {
		p1 := NewProperty(1)
		a,b := scene.Add("a"),scene.Add("b")
		a.Add(p1)
		b.Add(p1)

		c.Specify("both actors contain same property", func() {
			c.Expect(a.properties[1], IsSame, b.properties[1])
		})

		c.Specify("scene contains both actors", func() {
			c.Expect(len(scene.Actors), Equals, 2)
			c.Expect(scene.Actors["a"], IsSame, a)
			c.Expect(scene.Actors["b"], IsSame, b)
			c.Expect(len(scene.byProperty), Equals, 1)
			c.Expect(scene.byProperty[1], Contains, a)
			c.Expect(scene.byProperty[1], Contains, b)
		})

	})
}

func RemovePropertySpec(c gospec.Context) {
	scene := New()

	c.Specify("Actor with two properties", func() {
		p1, p2 := NewProperty(1), NewProperty(2)
		a := scene.Add("a")
		a.Add(p1)
		a.Add(p2)

		c.Specify("removing one", func() {
			a.Remove(1)

			c.Specify("actor has only the other property", func() {
				c.Expect(len(a.properties), Equals, 1)
				c.Expect(a.properties[2], IsSame, p2)
			})

		})
	})
}

func RemoveTypeSpec(c gospec.Context) {
	scene := New()

	c.Specify("Actor with two properties of different type", func() {
		p1, p2 := NewProperty(1), NewProperty(2)
		a := scene.Add("a")
		a.Add(p1)
		a.Add(p2)

		c.Specify("remove one of those properties", func() {
			r,present := a.Remove(1)

			c.Specify("present is true", func() {
				c.Expect(present, IsTrue)
			})
			
			c.Specify("actor has only other property left", func() {
				c.Expect(len(a.properties), Equals, 1)
				c.Expect(a.properties[2], IsSame, p2)
			})

			c.Specify("returns the removed property", func() {
				c.Expect(r, IsSame, p1)
			})
		})

		c.Specify("remove both property types", func() {
			a.Remove(1)
			a.Remove(2)

			c.Specify("actor has no properties left", func() {
				c.Expect(len(a.properties), Equals, 0)
			})
		})
	})

	c.Specify("Removing same actor twice", func() {
		p1 := NewProperty(1)
		a,b := scene.Add("a"),scene.Add("b")
		a.Add(p1)
		b.Add(p1)
		scene.Remove(a)
		scene.Remove(a)

		c.Specify("does not affect other actor", func() {
			c.Expect(len(scene.Actors), Equals, 1)
			c.Expect(scene.Actors["b"], IsSame, b)
			c.Expect(len(a.properties), Equals, 1)
			c.Expect(a.properties[1], IsSame, p1)
		})
	})
}

func FindSpec(c gospec.Context) {
	scene := New()

	c.Specify("Find on empty scene returns empty list", func() {
		result := scene.Find(1)
		c.Expect(len(result), Equals, 0)
	})

	c.Specify("One actor, one property", func() {
		p1 := NewProperty(1)
		a := scene.Add("a")
		a.Add(p1)

		c.Specify("for that property returns that actor", func() {
			result := scene.Find(1)
			c.Expect(len(result), Equals, 1)
			c.Expect(result, Contains, a)
		})

		c.Specify("for another property returns empty list", func() {
			result := scene.Find(2)
			c.Expect(len(result), Equals, 0)
		})
	})

	c.Specify("Two actors, sharing one property", func() {
		p1,p2 := NewProperty(1),NewProperty(2)
		a,b := scene.Add("a"),scene.Add("b")
		a.Add(p1)
		b.Add(p1)
		b.Add(p2)
		
		c.Specify("requesting shared property returns both", func() {
			result := scene.Find(1)

			c.Expect(len(result), Equals, 2)
			c.Expect(result, Contains, a)
			c.Expect(result, Contains, b)
		})

		c.Specify("requesting property specific to just one returns that one", func() {
			result := scene.Find(2)
			c.Expect(len(result), Equals, 1)
			c.Expect(result, Contains, b)
		})

		c.Specify("requesting property that none has returns empty list", func() {
			result := scene.Find(3)
			c.Expect(len(result), Equals, 0)
		})
	})
}