package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/Gurpartap/null/internal"
)

type JSONB struct {
	hasValue bool
	value    []byte
}

func NewJSONB(value []byte, hasValue bool) JSONB {
	opt := &JSONB{}
	if hasValue {
		opt.SetValue(value)
	}
	return *opt
}

// SetValue performs the conversion.
func (opt *JSONB) SetValue(value []byte) {
	opt.value = append(opt.value[0:0], value...)
	opt.hasValue = true
}

// Unwrap moves the value out of the optional, if it is Some(value).
// This function returns multiple values, and if that's undesirable,
// consider using Some and None functions.
func (opt JSONB) Unwrap() ([]byte, bool) {
	return opt.getValue(), opt.getHasValue()
}

// UnwrapOr returns the contained value or a default.
func (opt JSONB) UnwrapOr(def []byte) []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return def
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (opt JSONB) UnwrapOrElse(fn func() []byte) []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return fn()
}

// UnwrapOrDefault returns the contained value or the default.
func (opt JSONB) UnwrapOrDefault() []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	return nil
}

// UnwrapOrPanic returns the contained value or panics.
func (opt JSONB) UnwrapOrPanic() []byte {
	if opt.getHasValue() {
		return opt.getValue()
	}
	panic("unable to unwrap JSONB")
}

func (opt JSONB) getHasValue() bool {
	return opt.hasValue
}

func (opt JSONB) getValue() []byte {
	return opt.value
}

// String conforms to fmt Stringer interface.
func (opt JSONB) String() string {
	if value, ok := opt.Unwrap(); ok {
		return fmt.Sprintf("Some(%v)", value)
	}
	return "null"
}

// MarshalJSON implements the json Marshaler interface.
func (opt JSONB) MarshalJSON() ([]byte, error) {
	if !opt.getHasValue() {
		return []byte("null"), nil
	}
	return json.Marshal(opt.getValue())
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (opt *JSONB) UnmarshalJSON(data []byte) error {
	if data == nil || bytes.Equal(data, []byte("null")) {
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
func (opt *JSONB) Scan(src interface{}) error {
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
func (opt JSONB) Value() (driver.Value, error) {
	if !opt.getHasValue() || bytes.Equal(opt.getValue(), []byte("null")) {
		return nil, nil
	}

	if bytes.Equal(opt.getValue(), []byte{}) {
		return []byte("{}"), nil
	}

	// Remove NUL character byte(s)
	//
	// Postgres will reject jsonb with \u0000. For details, see
	// https://www.postgresql.org/docs/9.4/static/release-9-4-1.html
	//
	// "jsonb is stricter, and as such, it disallows Unicode escapes for
	// non-ASCII characters (those above U+007F) unless the database encoding
	// is UTF8. It also rejects the NULL character (\u0000), which cannot be
	// represented in PostgreSQL's text type."
	//
	// https://www.compose.com/articles/faster-operations-with-the-jsonb-data-type-in-postgresql/
	return bytes.Replace(opt.getValue(), []byte("\\u0000"), []byte{}, -1), nil
}
