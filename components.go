package components

import "time"

type Actor string
type Property uint

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
	Actors     map[Property][]Actor
	Properties map[Actor][]Property
}

func NewScene() *Scene {
	scene := Scene{}
	scene.Actors = map[Property][]Actor{}
	scene.Properties = map[Actor][]Property{}
	return &scene
}

func (s Scene) Add(a Actor, p Property) {
	if v, present := s.Properties[a]; !present {
		v = []Property{}
		s.Properties[a] = v
	}
	s.Properties[a] = append(s.Properties[a], p)

	if v, present := s.Actors[p]; !present {
		v = []Actor{}
		s.Actors[p] = v
	}
	s.Actors[p] = append(s.Actors[p], a)
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

func (s Scene) Find(p ...Property) []Actor {
	// return actors that have all of the provided properties
	result := []Actor{}

	if actors, present := s.Actors[p[0]]; present {
		for _, a := range actors {
			rest := p[1:]
			ap := s.Properties[a]

			// opt: if actor prop count < len(p) -> skip actor
			if len(ap) < len(p) {
				continue
			}
			all := true
			for _, wanted := range rest {
				this := false
				for _, exist := range ap {
					if exist == wanted {
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