package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Gurpartap/null/internal"
	"github.com/pkg/errors"
)

type String struct {
	hasValue bool
	value    string
}

func NewString(value string, hasValue bool) String {
	opt := &String{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *String) SetValue(value string) {
	opt.value = value
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt String) Unwrap() (string, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt String) UnwrapOr(def string) string {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt String) UnwrapOrElse(fn func() string) string {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt String) UnwrapOrDefault() string {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return ""
}

// UnwrapOrPanic returns the contained value or panics.
func (opt String) UnwrapOrPanic() string {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap String")
}

// Or returns the optional if it contains a value, otherwise returns optb.
func (opt String) Or(optb String) String {
	if opt.getHasValue() {
		return opt
	}
	return optb
}

func (opt String) getHasValue() bool {
	return opt.hasValue
}

func (opt String) getValue() string {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt String) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt String) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) || data == nil {
		opt.value, opt.hasValue = "", false
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
func (opt *String) Scan(src interface{}) error {
	if src == nil {
		opt.value, opt.hasValue = "", false
		return nil
	}

	var value string
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return err
	}
	opt.SetValue(value)

	return nil
}

// Value implements the driver Valuer interface.
func (opt String) Value() (driver.Value, error) {
	if !opt.getHasValue() {
		return nil, nil
	}
	return string(opt.getValue()), nil
}
