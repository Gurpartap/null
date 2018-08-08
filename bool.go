package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Gurpartap/null/internal"
	"github.com/pkg/errors"
)

type Bool struct {
	hasValue bool
	value    bool
}

func NewBool(value bool, hasValue bool) Bool {
	opt := &Bool{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *Bool) SetValue(value bool) {
	opt.value = value
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt Bool) Unwrap() (bool, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt Bool) UnwrapOr(def bool) bool {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt Bool) UnwrapOrElse(fn func() bool) bool {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt Bool) UnwrapOrDefault() bool {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return false
}

// UnwrapOrPanic returns the contained value or panics.
func (opt Bool) UnwrapOrPanic() bool {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap Int64")
}

func (opt Bool) getHasValue() bool {
	return opt.hasValue
}

func (opt Bool) getValue() bool {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt Bool) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt Bool) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) || data == nil {
		opt.value, opt.hasValue = false, false
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
func (opt *Bool) Scan(src interface{}) error {
	if src == nil {
		opt.value, opt.hasValue = false, false
		return nil
	}

	var value bool
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return err
	}
	opt.SetValue(value)

	return nil
}

// Value implements the driver Valuer interface.
func (opt Bool) Value() (driver.Value, error) {
	if !opt.getHasValue() {
		return nil, nil
	}
	return opt.getValue(), nil
}
