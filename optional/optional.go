package optional

// Optional is a type which may hold value or not.
type Optional[T any] struct {
	val   T
	isSet bool
}

func (o *Optional[T]) Set(v T) {
	o.val = v
	o.isSet = true
}

func (o *Optional[T]) Get() T {
	return o.val
}

func (o *Optional[T]) HasValue() bool {
	return o.isSet
}
