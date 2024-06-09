package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type NullableType interface {
	any | bool | string | float64 | int16 | int32 | int64 | time.Time
}
type Nullable[T any] struct {
	Value NullableType
	Set   bool
}

// If this method was called, the value was set.
func (i *Nullable[T]) UnmarshalJSON(data []byte) error {
	i.Set = true
	var temp NullableType
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	switch temp := temp.(type) {
	case string:
		i.Value = sql.NullString{String: temp, Valid: temp != ""}
	case bool:
		i.Value = sql.NullBool{Bool: temp, Valid: !temp}
	case float64:
		i.Value = sql.NullFloat64{Float64: temp, Valid: true}
	case int16:
		i.Value = sql.NullInt16{Int16: temp, Valid: true}
	case int32:
		i.Value = sql.NullInt32{Int32: temp, Valid: true}
	case int64:
		i.Value = sql.NullInt64{Int64: temp, Valid: true}
	case time.Time:
		i.Value = sql.NullTime{Time: temp, Valid: true}
	default:
		return fmt.Errorf("unsupported type %T for nullable field", temp)
	}

	return nil
}

func (i *Nullable[T]) ToNullableAny() Nullable[any] {
	value := any(i.Value)
	return Nullable[any]{Value: &value, Set: i.Set}
}
