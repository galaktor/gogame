/*
 Contains basic types and function required to manage a scene. 
*/
package scene

import "fmt"
import "errors"

// An actor has 0 or more Properties, which represent a type of data
// that the describes the actor.
type Actor struct {
	Id         string
	properties map[PropertyType]Property
	s          *Scene
}

func newActor(id string) *Actor {
	return &Actor{Id: id, properties: map[PropertyType]Property{}}

}

// Adds a property to the actor. Returns an error if the Actor already
// has a property of that type id. Actors can only hold one property
// of a given type.
func (a *Actor) Add(p Property) error {
	t := p.Type()
	if _, present := a.properties[t]; present {
		msg := fmt.Sprintf("actor %v already contains property of type %v", a.Id, t)
		return errors.New(msg)
	}

	a.properties[t] = p
	a.s.cache(a, t)

	return nil
}

// Retrieves a property of a given type from an actor.
// Returns nil if the actor does not have a property of
// the requested type.
func (a *Actor)Get(p PropertyType) Property {
	return a.properties[p]
}

// Removes a property from an actor and returns a pointer to the
// removed property. "present" will be false, and removed nil if
// the property type was not found on teh actor.
func (a *Actor)Remove(t PropertyType) (removed Property,present bool) {
	removed,present = a.properties[t]
	delete(a.properties, t)
	a.s.uncache(a, t)
	return
}

// TODO: use actual types as Id instead
// get type without instance:
// reflect.TypeOf((*MyType)(nil)).Elem()
// -> cast nil to wanted type, use reflect Elem() to get type
type PropertyType uint

// An interface that all properties in the scene need to implement.
type Property interface {
	Type() PropertyType
}

// A scene has Actors, identifiable by a string (unique per scene).
// The scene also indexes actors by their property types for faster
// lookup in Find()
// Properties should hold /only/ data, and some functions related 
// to managing that data, no game or system related logic.
type Scene struct {
	byProperty map[PropertyType][]*Actor
	Actors     map[string]*Actor
}

// Creates a new, empty scene.
func New() *Scene {
	return &Scene{map[PropertyType][]*Actor{}, map[string]*Actor{}}
}

// Adds an actor with the given id to the scene and returns a pointer to it.
// Will panic if the scene already contains an actor with that id.
// Actor contains a reference to the scene that created it in order to
// keep the scene's property index up-to-date when properties are added or
// removed on the actor.
func (s *Scene) Add(id string) *Actor {
	a := newActor(id)
	e := s.addActor(a)
	a.s = s

	if e != nil {
		panic(e.Error())
	}

	return a
	
}

func (s Scene) addActor(a *Actor) error {
	if _,present := s.Actors[a.Id]; present {
		msg := fmt.Sprintf("scene already contains actor with id %v", a.Id)
		return errors.New(msg)
	}

	s.Actors[a.Id] = a
	for t := range a.properties {
		s.cache(a, t)
	}

	return nil
}

// Removes a given actor from the scene.
func (s Scene) Remove(a *Actor) {
	if _, present := s.Actors[a.Id]; !present {
		return
	}

	delete(s.Actors, a.Id)

	for t := range a.properties {
		s.uncache(a, t)
	}
}

func (s Scene)cache(a *Actor, t PropertyType) {
	if _,present := s.byProperty[t]; !present {
		s.byProperty[t] = []*Actor{}
	}
	s.byProperty[t] = append(s.byProperty[t], a)
}

func (s Scene)uncache(a *Actor, t PropertyType) {
	if actors,present := s.byProperty[t]; present {
		// TODO: pre-allocate right size rather than constant resizing
		newlist := []*Actor{}
		for _,actor := range actors {
			// keep all but the uncached one
			if actor.Id != a.Id {
				newlist = append(newlist,actor)
			}
		}
		s.byProperty[t] = newlist
	}
}

// Allows for very specialized query of the sccene by property type.
// Given one or more property types, will return a list of all actors
// that contain every given type. For very large scenes this will
// probably have to be improved in many ways, possibly by using a
// binary search tree.
func (s Scene) Find(p ...PropertyType) (result []*Actor) {
	// opt: exclude actors without first property
	if actors, present := s.byProperty[p[0]]; present {
		if len(p) == 1 {
			// opt: quit now if only looking for one property
			result = actors
			return
		}

		for _, a := range actors {
			if len(p) > len(a.properties) {
				// opt: exclude actors with less properties than requested
				// this assumes that actors only have one of each property type
				continue
			}

			// opt: we already checked prop at 0
			rest := p[1:] 
			// opt: requested less or equal props than actor has
			// loop on those rather than all the props of the actor
			// look until we find a prop that doesn't match
			hit := true
			for _, wanted := range rest {
				if _,present := a.properties[wanted]; !present {
					hit = false
					break
				}
			}

			if hit {
				result = append(result, a)
			}
		}
	}

	return result
}
