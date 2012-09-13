package scene

// TODO: use actual types as Id instead?
// get type without instance:
// reflect.TypeOf((*MyType)(nil)).Elem()
// -> cast nil to wanted type, use reflect Elem() to get type
type PropertyType uint

// An interface that all properties in the scene need to implement.
type Property interface {
	Type() PropertyType
}