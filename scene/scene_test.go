package scene

import (
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
)

func CtorSpec(c gospec.Context) {
	c.Specify("New scene", func() {
		scene := NewScene()

		c.Specify("has zero actors", func() {
			c.Expect(len(scene.Actors), Equals, 0)
	})

		c.Specify("has zero properties", func() {
			c.Expect(len(scene.Properties), Equals, 0)
		})
	})
}

type SomeProperty struct {
	tid PropertyType
}

func NewProperty(t PropertyType) Property {
	return &SomeProperty{t}
}

func (p *SomeProperty)Type() PropertyType {
	return p.tid
}


func AddSpec(c gospec.Context) {
	scene := NewScene()
	
	c.Specify("Add one actor with one property", func() {
		p1 := NewProperty(1)
		scene.Add("a", p1)

		c.Specify("scene contains one property", func() {
			c.Expect(len(scene.Properties), Equals, 1)
			c.Expect(len(scene.Properties["a"]), Equals, 1)
		})

		c.Specify("scene contains just that property for that actor", func() {
			c.Expect(scene.Properties["a"][0], Equals, p1)
		})

		c.Specify("scene contains one actor", func() {
			c.Expect(len(scene.Actors), Equals, 1)
			c.Expect(len(scene.Actors[1]), Equals, 1)
		})

		c.Specify("scene contains just that actor for that property", func() {
			c.Expect(scene.Actors[1][0], Equals, Actor("a"))
		})
	})

	c.Specify("Add two actors with one different property each", func() {
		p1 := NewProperty(1)
		p2 := NewProperty(2)
		scene.Add("a", p1)
		scene.Add("b", p2)
		
		c.Specify("scene contains both properties", func() {
			c.Expect(len(scene.Properties), Equals, 2)
			c.Expect(len(scene.Properties["a"]), Equals, 1)
			c.Expect(scene.Properties["a"][0], Equals, p1)
			c.Expect(len(scene.Properties["b"]), Equals, 1)
			c.Expect(scene.Properties["b"][0], Equals, p2)
		})

		c.Specify("scene contains both actors", func() {
			c.Expect(len(scene.Actors), Equals, 2)
			c.Expect(len(scene.Actors[1]), Equals, 1)
			c.Expect(scene.Actors[1][0], Equals, Actor("a"))
			c.Expect(len(scene.Actors[2]), Equals, 1)
			c.Expect(scene.Actors[2][0], Equals, Actor("b"))
		})
	})

	c.Specify("Add two properties to same actor", func() {
		p1 := NewProperty(1)
		p2 := NewProperty(2)
		scene.Add("a", p1)
		scene.Add("a", p2)

		c.Specify("scene contains both properties", func() {
			c.Expect(len(scene.Properties["a"]), Equals, 2)
			c.Expect(scene.Properties["a"], Contains, p1)
			c.Expect(scene.Properties["a"], Contains, p2)
		})
	})

	c.Specify("Add two actors with same property", func() {
		p1 := NewProperty(1)
		scene.Add("a", p1)
		scene.Add("b", p1)

		c.Specify("scene contains both actors", func() {
			c.Expect(len(scene.Actors[1]), Equals, 2)
			c.Expect(scene.Actors[1], Contains, Actor("a"))
			c.Expect(scene.Actors[1], Contains, Actor("b"))
		})

		c.Specify("scene contains just that property", func() {
			c.Expect(len(scene.Properties["a"]), Equals, 1)
			c.Expect(scene.Properties["a"], Contains, p1)
			c.Expect(len(scene.Properties["b"]), Equals, 1)
			c.Expect(scene.Properties["b"], Contains, p1)
		})
	})
}

func RemovePropertySpec(c gospec.Context) {
	scene := NewScene()

	c.Specify("Actor with two properties of same type", func() {
		p1, p2 := NewProperty(1), NewProperty(2)
		scene.Add("a", p1)
		scene.Add("a", p2)

		c.Specify("removing one", func() {
			scene.RemoveProperty("a", p1)

			c.Specify("actor has only the other property", func() {
				c.Expect(len(scene.Properties["a"]), Equals, 1)
				c.Expect(scene.Properties["a"], Contains, p2)
			})

			
		})
	})
}

func RemoveTypeSpec(c gospec.Context) {
	scene := NewScene()

	c.Specify("Actor with two properties of different type", func() {
		p1, p2 := NewProperty(1), NewProperty(2)
		scene.Add("a", p1)
		scene.Add("a", p2)

		c.Specify("remove one property type", func() {
			ret := scene.RemoveType("a", p2.Type())
			
			c.Specify("actor has only other property left", func() {
				c.Expect(len(scene.Properties["a"]), Equals, 1)
				c.Expect(scene.Properties["a"], Contains, p1)
			})

			c.Specify("returns the removed property", func() {
				c.Expect(len(ret), Equals, 1)
				c.Expect(ret, Contains, p2)
			})
		})

		c.Specify("remove both property types", func() {
			scene.RemoveType("a", p1.Type())
			scene.RemoveType("a", p2.Type())

			c.Specify("actor has no properties left", func() {
				c.Expect(len(scene.Properties["a"]), Equals, 0)
			})
		})
	})

	c.Specify("Removing unknown actor does not affect other actor", func() {
		scene.Add("a", NewProperty(1))
		ret := scene.RemoveType("b", 1)

		c.Specify("does not affect other actor", func() {
			c.Expect(len(scene.Properties["a"]), Equals, 1)
			c.Expect(len(scene.Actors[1]), Equals, 1)
		})
		
		c.Specify("returns empty list", func() {
			c.Expect(len(ret), Equals, 0)
		})
		
		
	})
	
	c.Specify("Actor with two properties of same type", func() {
		p1, p2 := NewProperty(1), NewProperty(1)
		scene.Add("a", p1)
		scene.Add("a", p2)

		c.Specify("removing that property type", func() {
			ret := scene.RemoveType("a", 1)

			c.Specify("leaves the actor with zero properties", func() {
				c.Expect(len(scene.Properties["a"]), Equals, 0)
			})

			c.Specify("returns both removed properties", func() {
				c.Expect(len(ret), Equals, 2)
				c.Expect(ret, Contains, p1)
				c.Expect(ret, Contains, p2)
			})
		})
		
		
	})
}

func FindSpec(c gospec.Context) {
	scene := NewScene()
	
	c.Specify("Find on empty scene returns empty list", func() {
		result := scene.Find(1)
		c.Expect(len(result), Equals, 0)
	})

	c.Specify("One actor, one property", func() {
		p1 := NewProperty(1)
		scene.Add("a", p1)

		c.Specify("for that property returns that actor", func() {
			result := scene.Find(1)
			c.Expect(len(result), Equals, 1)
			c.Expect(result[0], Equals, Actor("a"))
		})

		c.Specify("for another property returns empty list", func() {
			result := scene.Find(2)
			c.Expect(len(result), Equals, 0)
		})
	})

	c.Specify("Two actors, sharing one property", func() {
		p1 := NewProperty(1)
		p2 := NewProperty(2)
		scene.Add("a", p1)
		scene.Add("b", p1)
		scene.Add("b", p2)

		
		c.Specify("requesting shared property returns both", func() {
			result := scene.Find(1)
			c.Expect(len(result), Equals, 2)
		})
	
		c.Specify("requesting property specific to just one returns that one", func() {
			result := scene.Find(2)
			c.Expect(len(result), Equals, 1)
			c.Expect(result[0], Equals, Actor("b"))
		})

		c.Specify("requesting property that none has returns empty list", func() {
			result := scene.Find(3)
			c.Expect(len(result), Equals, 0)
		})
	})
}