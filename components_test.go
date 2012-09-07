package components

import (
	"github.com/orfjackal/gospec"
	. "github.com/orfjackal/gospec"
)

func SceneCtorSpec(c gospec.Context) {
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

func SceneAddSpec(c gospec.Context) {
	scene := NewScene()
	
	c.Specify("Add one actor with one property", func() {
		scene.Add("a", 1)

		c.Specify("scene contains one property", func() {
			c.Expect(len(scene.Properties), Equals, 1)
			c.Expect(len(scene.Properties["a"]), Equals, 1)
		})

		c.Specify("scene contains just that property for that actor", func() {
			c.Expect(scene.Properties["a"][0], Equals, Property(1))
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
		scene.Add("a", 1)
		scene.Add("b", 2)
		
		c.Specify("scene contains both properties", func() {
			c.Expect(len(scene.Properties), Equals, 2)
			c.Expect(len(scene.Properties["a"]), Equals, 1)
			c.Expect(scene.Properties["a"][0], Equals, Property(1))
			c.Expect(len(scene.Properties["b"]), Equals, 1)
			c.Expect(scene.Properties["b"][0], Equals, Property(2))
		})

		c.Specify("scene contains both actors", func() {
			c.Expect(len(scene.Actors), Equals, 2)
			c.Expect(len(scene.Actors[1]), Equals, 1)
			c.Expect(scene.Actors[1][0], Equals, Actor("a"))
			c.Expect(len(scene.Actors[2]), Equals, 1)
			c.Expect(scene.Actors[2][0], Equals, Actor("b"))
		})
	})

	c.Specify("Add two actors with same property", func() {
		scene.Add("a", 1)
		scene.Add("b", 1)

		c.Specify("scene contains both actors", func() {
			c.Expect(len(scene.Actors[1]), Equals, 2)
			c.Expect(scene.Actors[1], Contains, Actor("a"))
			c.Expect(scene.Actors[1], Contains, Actor("b"))
		})

		c.Specify("scene contains just that property", func() {
			c.Expect(len(scene.Properties["a"]), Equals, 1)
			c.Expect(scene.Properties["a"], Contains, Property(1))
			c.Expect(len(scene.Properties["b"]), Equals, 1)
			c.Expect(scene.Properties["b"], Contains, Property(1))
		})
	})
}

func SceneRemoveSpec(c gospec.Context) {
	c.Specify("I need tests!", func() {
		c.Expect(false, Equals, true)
	})
}

func SceneFindSpec(c gospec.Context) {
	scene := NewScene()
	
	c.Specify("Find on empty scene returns empty list", func() {
		result := scene.Find(1)
		c.Expect(len(result), Equals, 0)
	})

	c.Specify("One actor, one property", func() {
		scene.Actors = map[Property][]Actor{1:[]Actor{"a"}}
		scene.Properties = map[Actor][]Property{"a":[]Property{1}}

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
		scene.Actors = map[Property][]Actor{1:[]Actor{"a","b"},2:[]Actor{"b"}}
		scene.Properties = map[Actor][]Property{"a":[]Property{1},"b":[]Property{1,2}}
		
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