package scene

import "time"

type Actor string

type PropertyType uint

type Property interface {
	Type() PropertyType
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

// todo: custom, optimized data structure, e.g. binary search tree
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

func (s Scene) Add(a Actor, p Property) {
	if _, present := s.Properties[a]; !present {
		s.Properties[a] = []Property{}
	}
	s.Properties[a] = append(s.Properties[a], p)

	if _, present := s.Actors[p.Type()]; !present {
		s.Actors[p.Type()] = []Actor{}
	}
	s.Actors[p.Type()] = append(s.Actors[p.Type()], a)
}

func (s Scene)RemoveType(a Actor, t PropertyType) (removed []Property) {
	selector := func(p Property) bool {
		return p.Type() == t
	}

	removed = s.RemoveWhere(a, selector)
	return
}

func (s Scene)RemoveProperty(a Actor, toRemove Property) {
	selector := func(p Property) bool {
		return p == toRemove
	}

	_ = s.RemoveWhere(a, selector)
	return
}

type PropertySelector func(Property) bool

func (s Scene)RemoveWhere(a Actor, sel PropertySelector) (removed []Property) {	
	if props,present := s.Properties[a]; present {
		kept := []Property{}
		for _,prop := range props {
			if sel(prop) {
				removed = append(removed,prop)
			} else {
				kept = append(kept,prop)
			}
		}
		s.Properties[a] = kept
	}

	return removed
}




/*
func (s Scene)Remove(a Actor) {
	//delete(s.Properties, a)

	// remove actor from all property slices
}
*/

func (s Scene) Find(p ...PropertyType) (result []Actor) {
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
					if exist.Type() == wanted {
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