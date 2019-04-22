package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/Gurpartap/null/internal"
)

// Int64Slice is a sql scanner interface for using []int64 as postgres nullable arrays.
type Int64Slice struct {
	hasValue bool
	value    []int64
}

func NewInt64Slice(value []int64, hasValue bool) Int64Slice {
	opt := &Int64Slice{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *Int64Slice) SetValue(value []int64) {
	opt.value = append(opt.value[0:0], value...)
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt Int64Slice) Unwrap() ([]int64, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt Int64Slice) UnwrapOr(def []int64) []int64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt Int64Slice) UnwrapOrElse(fn func() []int64) []int64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt Int64Slice) UnwrapOrDefault() []int64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return nil
}

// UnwrapOrPanic returns the contained value or panics.
func (opt Int64Slice) UnwrapOrPanic() []int64 {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap Int64Slice")
}

func (opt Int64Slice) getHasValue() bool {
	return opt.hasValue
}

func (opt Int64Slice) getValue() []int64 {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt Int64Slice) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt Int64Slice) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *Int64Slice) UnmarshalJSON(data []byte) error {
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
func (opt *Int64Slice) Scan(src interface{}) error {
	if src == nil {
		opt.value, opt.hasValue = nil, false
		return nil
	}

	var value string
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return errors.WithStack(err)
	}

	slice, err := internal.StringToInt64Slice(value)
	if err != nil {
		return errors.WithStack(err)
	}

	opt.SetValue(slice)

	return nil
}

// Value implements the driver Valuer interface.
func (opt Int64Slice) Value() (driver.Value, error) {
	if !opt.getHasValue() {
		return nil, nil
	}

	var s []string
	for _, i := range opt.getValue() {
		s = append(s, strconv.FormatInt(i, 10))
	}
	return "{" + strings.Join(s, ",") + "}", nil
}
