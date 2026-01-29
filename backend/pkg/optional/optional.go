package optional

type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional with a value.
func Some[T any](v T) Optional[T] { return Optional[T]{value: v, present: true} }

// None creates an Optional without a value.
func None[T any]() Optional[T] {
	var zero T
	return Optional[T]{value: zero, present: false}
}

// IsSome returns true if value is present.
func (o Optional[T]) IsSome() bool { return o.present }

// IsNone returns true if no value.
func (o Optional[T]) IsNone() bool { return !o.present }

// Unwrap returns the value (zero value if None).
func (o Optional[T]) Unwrap() T { return o.value }

// UnwrapOr returns the value or a fallback.
func (o Optional[T]) UnwrapOr(fallback T) T {
	if o.present {
		return o.value
	}
	return fallback
}

// UnwrapOrAny returns the value or a fallback of any type.
func (o Optional[T]) UnwrapOrAny(fallback any) any {
	if o.present {
		return o.value
	}
	return fallback
}
