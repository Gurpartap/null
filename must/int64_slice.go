package must

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/Gurpartap/null/internal"
)

// Int64Slice is a sql scanner interface for using []int64 as postgres arrays.
type Int64Slice []int64

// MarshalJSON implements the json Marshaler interface.
func (v Int64Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (v *Int64Slice) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Scan implements the sql Scanner interface.
func (v *Int64Slice) Scan(src interface{}) error {
	var value string
	err := internal.ConvertAssign(&value, src)
	if err != nil {
		return errors.WithStack(err)
	}

	slice, err := internal.StringToInt64Slice(value)
	if err != nil {
		return errors.WithStack(err)
	}

	*v = append((*v)[0:0], slice...)

	return nil
}

// Value implements the driver Valuer interface.
func (v Int64Slice) Value() (driver.Value, error) {
	var s []string
	for _, i := range v {
		s = append(s, strconv.FormatInt(i, 10))
	}
	return "{" + strings.Join(s, ",") + "}", nil
}
