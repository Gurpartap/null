package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/Gurpartap/null/internal"
)

type Bytes struct {
	hasValue bool
	value    []byte
}

func NewBytes(value []byte, hasValue bool) Bytes {
	opt := &Bytes{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *Bytes) SetValue(value []byte) {
	opt.value = value
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt Bytes) Unwrap() ([]byte, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt Bytes) UnwrapOr(def []byte) []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt Bytes) UnwrapOrElse(fn func() []byte) []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt Bytes) UnwrapOrDefault() []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return nil
}

// UnwrapOrPanic returns the contained value or panics.
func (opt Bytes) UnwrapOrPanic() []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap Bytes")
}

func (opt Bytes) getHasValue() bool {
	return opt.hasValue
}

func (opt Bytes) getValue() []byte {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt Bytes) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt Bytes) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *Bytes) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) || data == nil {
		opt.value, opt.hasValue = nil, false
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
func (opt *Bytes) Scan(src interface{}) error {
	if src == nil {
		opt.value, opt.hasValue = nil, false
		return nil
	}

	var value []byte
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return errors.WithStack(err)
	}
	opt.SetValue(value)

	return nil
}

// Value implements the driver Valuer interface.
func (opt Bytes) Value() (driver.Value, error) {
	if !opt.getHasValue() {
		return nil, nil
	}
	return []byte(opt.getValue()), nil
}
