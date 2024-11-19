package errors

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// New returns an error that formats as the given code and text (the message).
// Each call to New returns a distinct error value even if the text is identical.
func New(code string, text string) error {
	return &Error{code, text}
}

// errorString is a trivial implementation of error.
type Error struct {
	Code    string
	Message string
}

// Returns the Web3Error code
//
// If e is not type of Web3Error returns an empty string
func Code(e error) string {
	if t, ok := e.(*Error); ok {
		return t.Code
	}
	return ""
}

// Returns the Web3Error message
//
// If e is not type of Web3Error returns an empty string
func Message(e error) string {
	if t, ok := e.(*Error); ok {
		return t.Message
	}
	return ""
}
func (e Error) Error() string {
	return e.Code + " - " + e.Message
}

func GetWeb3Error(err error) *Error {
	if w, ok := err.(*Error); ok {
		return w
	}
	return &Error{Message: err.Error()}
}

func (e *Error) String() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

// Implements the driver.Valuer interface to be able to insert it into MySQL
func (e Error) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Implements the sql.Scanner interface to be able to read it from MySQL
func (e *Error) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return New("ERROR_EXPECTS_BYTE_ARRAY", "expects []byte")
	}
	return json.Unmarshal(b, e)
}

// Implements the json.Marshaler interface to control JSON serialization.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.Code,
		Message: e.Message,
	})
}
