package optional

import (
	"encoding/json"
)

func (o Optional[T]) IsZero() bool { return o.IsNone() }

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.present {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	// Missing field or explicit null → None
	if string(data) == "null" {
		o.present = false
		var zero T
		o.value = zero
		return nil
	}

	// Otherwise try unmarshal into T
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.value = v
	o.present = true
	return nil
}
