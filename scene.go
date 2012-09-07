package scene

import "time"

type Actor string
type PropertyType uint

type Property struct {
	Tid PropertyType
}

type System interface {
	Update(timestep time.Duration)
}

func Start(s System, interval time.Duration) {
	ticker := time.Tick(interval)
	last := time.Now()
	for now := range ticker {
		s.Update(now.Sub(last))
		last = now
	}
}

// todo: custom, optimized data structure, i.e. search tree
type Scene struct {
	Actors     map[PropertyType][]Actor
	Properties map[Actor][]Property
}

func NewScene() *Scene {
	scene := &Scene{}
	scene.Actors = map[PropertyType][]Actor{}
	scene.Properties = map[Actor][]Property{}
	return scene
}

// do not take pointers; guarantee that every actor has their own copy
func (s Scene) Add(a Actor, p Property) {
	if _, present := s.Properties[a]; !present {
		s.Properties[a] = []Property{}
	}
	s.Properties[a] = append(s.Properties[a], p)

	if _, present := s.Actors[p.Tid]; !present {
		s.Actors[p.Tid] = []Actor{}
	}
	s.Actors[p.Tid] = append(s.Actors[p.Tid], a)
}

/*
func (s Scene)Remove(a Actor, p Property) {	
	// get actor, if present
	// find prop, if present (just loop for now)
	// re-slice around index of that prop

	// remove that actor from all property slices
}

func (s Scene)Remove(a Actor) {
	//delete(s.Properties, a)

	// remove actor from all property slices
}
*/

func (s Scene) Find(p ...PropertyType) []Actor {
	// return actors that have all of the provided properties
	result := []Actor{}

	if actors, present := s.Actors[p[0]]; present {
		for _, a := range actors {
			rest := p[1:]
			ap := s.Properties[a]

			// opt: if actor prop count < len(p) -> skip actor
			// WARNING: this assumes that actors only have one of each property type
			if len(ap) < len(p) {
				continue
			}

			all := true
			for _, wanted := range rest {
				this := false
				for _, exist := range ap {
					if exist.Tid == wanted {
						this = true
					}
				}
				if !this {
					all = false
					break
				}
			}
		
	
			if all {
				result = append(result,a)
			}
		}
	}

	return result
}