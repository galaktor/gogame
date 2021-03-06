package scene

import (
	"errors"
	"fmt"
)

// An actor has 0 or more Properties, which represent a type of data
// that the describes the actor.
type A struct {
	Id         string
	properties map[PType]P
	s          *S
}

func newActor(id string) *A {
	return &A{Id: id, properties: map[PType]P{}}

}

// Adds a property to the actor. Returns an error if the Actor already
// has a property of that type id. Actors can only hold one property
// of a given type.
func (a *A) Add(p P) error {
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
func (a *A)Get(p PType) P {
	return a.properties[p]
}

// Removes a property from an actor and returns a pointer to the
// removed property. "present" will be false, and removed nil if
// the property type was not found on teh actor.
func (a *A)Remove(t PType) (removed P,present bool) {
	removed,present = a.properties[t]
	delete(a.properties, t)
	a.s.uncache(a, t)
	return
}
