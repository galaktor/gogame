/*
 Contains basic types and function required to manage a scene. 
*/
package scene

import (
	"fmt"
	"errors"
)


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
