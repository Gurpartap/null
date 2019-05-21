# Nullable types for Go

```sh
go get -u github.com/Gurpartap/null
```

##### null
[![GoDoc](https://godoc.org/github.com/Gurpartap/null?status.svg)](https://godoc.org/github.com/Gurpartap/null)

##### must (non-nullable JSONB and Int64Slice)
[![GoDoc](https://godoc.org/github.com/Gurpartap/null/must?status.svg)](https://godoc.org/github.com/Gurpartap/null/must)

### Usage

```sh
go get -u github.com/Gurpartap/null
```
```go
import "github.com/Gurpartap/null"
```

#### Creating new nullable
```go
user := struct {
	Name null.String `db:"name" json:"name"`
}{}

user.Name = null.NewString("Badam Rogan", true)
```

#### Scaning from database
```go
_ = sqlxDB.Get(&user, `select name from users limit 1`)
```

#### Parsing from JSON
```go
b := []byte(`{"user": "Badam Rogan"}`)
_ = json.Unmarshal(b, &user)
```

#### Reading value
```go
fmt.Println(user.Name) // => Some("Badam Rogan")

if name, ok := user.Name.Unwrap(); ok {
	// has value
	fmt.Println(name) // => Badam Rogan
} else {
	// has no value
}
```

### Available methods

```go
func (opt null.String) Or(optb null.String) null.String {}
func (opt *null.String) SetValue(value string) {}
func (opt null.String) Unwrap() (string, bool) {}
func (opt null.String) UnwrapOr(def string) string {}
func (opt null.String) UnwrapOrDefault() string {}
func (opt null.String) UnwrapOrElse(fn func() string) string {}
func (opt null.String) UnwrapOrPanic() string {}
```

See godocs for full list and comments

## Credits

This package is inspired by, and inherits bits of code and experience from each of the following:

- https://github.com/piotrkowalczuk/ntypes
- https://github.com/leighmcculloch/go-optional
- https://github.com/markphelps/optional
- https://doc.rust-lang.org/std/option/
- https://doc.rust-lang.org/std/option/enum.Option.html
- https://github.com/apple/swift/blob/master/stdlib/public/core/Optional.swift
- https://github.com/Gurpartap/safer-go

## License

Copyright 2017 Gurpartap Singh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
