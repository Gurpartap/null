package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Gurpartap/null/internal"
	"github.com/pkg/errors"
)

type Time struct {
	hasValue bool
	value    time.Time
}

func NewTime(value time.Time, hasValue bool) Time {
	opt := &Time{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *Time) SetValue(value time.Time) {
	opt.value = value
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt Time) Unwrap() (time.Time, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt Time) UnwrapOr(def time.Time) time.Time {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt Time) UnwrapOrElse(fn func() time.Time) time.Time {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt Time) UnwrapOrDefault() time.Time {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return time.Time{}
}

// UnwrapOrPanic returns the contained value or panics.
func (opt Time) UnwrapOrPanic() time.Time {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap Time")
}

func (opt Time) getHasValue() bool {
	return opt.hasValue
}

func (opt Time) getValue() time.Time {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt Time) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt Time) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) || data == nil {
		opt.value, opt.hasValue = time.Time{}, false
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
func (opt *Time) Scan(src interface{}) error {
	if src == nil {
		opt.value, opt.hasValue = time.Time{}, false
		return nil
	}

	var value time.Time
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return err
	}
	opt.SetValue(value)

	return nil
}

// Value implements the driver Valuer interface.
func (opt Time) Value() (driver.Value, error) {
	if !opt.getHasValue() {
		return nil, nil
	}
	return time.Time(opt.getValue()), nil
}
