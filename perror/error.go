package perror

import (
	"fmt"
)

const (
	ErrTypeWarn   = "warn"
	ErrTypeError  = "err"
	ErrTypeIgnore = "ignore"
)

// NewPError create new error with status and message
// errType=warn
func New(status int, msg ...interface{}) *Error {

	return initError(ErrTypeError, status, msg...)
}

// NewWarn create new error with status and message
// errType=warn
func NewWarn(status int, msg ...interface{}) *Error {
	return initError(ErrTypeWarn, status, msg...)
}

// NewIgnore create new error with status and message
// errType=ignore
func NewIgnore(status int, msg ...interface{}) *Error {
	return initError(ErrTypeIgnore, status, msg...)
}

func initError(errType string, status int, msg ...interface{}) *Error {
	message := ""
	if len(msg) == 1 {
		message = msg[0].(string)
	} else if len(msg) > 1 {
		message = fmt.Sprintf(msg[0].(string), msg[1:]...)
	}

	return &Error{errType, status, message}
}

// Exception panic as exception
type Error struct {
	errType string
	status  int
	message string
}

// ErrType return errType
func (p *Error) ErrType() string {
	return p.errType
}

// IsErrTypeWarn return is warn
func (p *Error) IsErrTypeWarn() bool {
	return p.errType == ErrTypeWarn
}

// IsErrTypeErr return is err
func (p *Error) IsErrTypeErr() bool {
	return p.errType == ErrTypeError
}

// IsErrTypeIgnore return is ignore
func (p *Error) IsErrTypeIgnore() bool {
	return p.errType == ErrTypeIgnore
}

// GetStatus get exception status code
func (p *Error) Status() int {
	return p.status
}

// GetMessage get exception message string
func (p *Error) Message() string {
	return p.message
}

// Error implement error interface
func (p *Error) Error() string {
	return fmt.Sprintf("errCode: %d, errMsg: %s", p.status, p.message)
}
