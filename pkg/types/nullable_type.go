package types

import "encoding/json"

type Nullable[T any] struct {
	Value *T
	Set   bool
}

// If this method was called, the value was set.
func (i *Nullable[T]) UnmarshalJSON(data []byte) error {
	i.Set = true
	var temp *T
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	return nil
}

func (i *Nullable[T]) ToNullableAny() Nullable[any] {
	if i.Value == nil {
		return Nullable[any]{Value: nil, Set: i.Set}
	}
	value := any(*i.Value)
	return Nullable[any]{Value: &value, Set: i.Set}
}
