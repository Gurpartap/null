package must

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/Gurpartap/null/internal"
)

type JSONB []byte

// MarshalJSON implements the json Marshaler interface.
func (v JSONB) MarshalJSON() ([]byte, error) {
	if bytes.Equal(v, []byte{}) || bytes.Equal(v, []byte("null")) {
		return []byte("{}"), nil
	}

	return json.Marshal(v)
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (v *JSONB) UnmarshalJSON(data []byte) error {
	if data == nil || bytes.Equal(data, []byte("null")) {
		*v = []byte{}
		return nil
	}

	var value []byte
	err := json.Unmarshal(data, &value)
	if err != nil {
		return errors.WithStack(err)
	}

	*v = append((*v)[0:0], value...)

	return nil
}

// Scan implements the sql Scanner interface.
func (v *JSONB) Scan(src interface{}) error {
	if src == nil {
		*v = []byte{}
		return nil
	}

	var value []byte
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return err
	}

	*v = append((*v)[0:0], value...)

	return nil
}

func (v JSONB) Value() (driver.Value, error) {
	if bytes.Equal(v, []byte{}) || bytes.Equal(v, []byte("null")) {
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
	return bytes.Replace(v, []byte("\\u0000"), []byte{}, -1), nil
}
