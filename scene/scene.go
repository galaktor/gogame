package scene

import "fmt"
import "errors"

type ActorId string

type Actor struct {
	Id         ActorId
	properties map[PropertyType]Property
	s          *Scene
}

func newActor(id ActorId) *Actor {
	return &Actor{Id: id, properties: map[PropertyType]Property{}}

}

func (a *Actor) Add(p Property) error {
	t := p.Type()
	if v, present := a.properties[t]; present {
		msg := fmt.Sprintf("actor %v already contains property of type %v", a.Id, t)
		return errors.New(msg)
	}

	a.properties[t] = p
	a.s.cache(a, t)

	return nil
}

func (a *Actor)Get(p PropertyType) (Property, bool) {
	prop,present := a.properties[p]
	return prop,present
}

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

type Property interface {
	Type() PropertyType
}

type Scene struct {
	byProperty map[PropertyType][]*Actor
	Actors     map[ActorId]*Actor
}

func NewScene() *Scene {
	return &Scene{map[PropertyType][]*Actor{}, map[ActorId]*Actor{}}
}

func (s Scene) Add(id ActorId) (a *Actor, e error) {
	a = newActor(ActorId(id))
	e = s.addActor(a)
	return
	
}

func (s Scene) addActor(a *Actor) error {
	if v,present := s.Actors[a.Id]; present {
		msg := fmt.Sprintf("scene already contains actor with id %v", a.Id)
		return errors.New(msg)
	}

	s.Actors[a.Id] = a
	for t := range a.properties {
		s.cache(a, t)
	}

	return nil
}

func (s Scene) Remove(a *Actor) {
	if actor, present := s.Actors[a.Id]; !present {
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
		newlist := make([]*Actor,len(actors)-1)
		for i,actor := range actors {
			// keep all but the uncached one
			if actor.Id != a.Id {
				newlist[i] = actor
			}
		}
		s.byProperty[t] = newlist
	}
}

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
