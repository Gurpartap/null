package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Gurpartap/null/internal"
	"github.com/pkg/errors"
)

type Float64 struct {
	hasValue bool
	value    float64
}

func NewFloat64(value float64, hasValue bool) Float64 {
	opt := &Float64{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *Float64) SetValue(value float64) {
	opt.value = value
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt Float64) Unwrap() (float64, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt Float64) UnwrapOr(def float64) float64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt Float64) UnwrapOrElse(fn func() float64) float64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt Float64) UnwrapOrDefault() float64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return 0
}

// UnwrapOrPanic returns the contained value or panics.
func (opt Float64) UnwrapOrPanic() float64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap Float64")
}

func (opt Float64) getHasValue() bool {
	return opt.hasValue
}

func (opt Float64) getValue() float64 {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt Float64) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt Float64) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *Float64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) || data == nil {
		opt.value, opt.hasValue = 0, false
		return nil
	}

	err := json.Unmarshal(data, &opt.value)
	if err != nil {
		opt.hasValue = false
		return errors.WithStack(err)
	}
	opt.hasValue = true

	return nil
}

// Scan implements the sql Scanner interface.
func (opt *Float64) Scan(src interface{}) error {
	if src == nil {
		opt.value, opt.hasValue = 0, false
		return nil
	}

	var value float64
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return err
	}
	opt.SetValue(value)

	return nil
}

// Value implements the driver Valuer interface.
func (opt Float64) Value() (driver.Value, error) {
	if !opt.getHasValue() {
		return nil, nil
	}
	return float64(opt.getValue()), nil
}
