package perror

import (
    "fmt"
)

// NewPError create new error with status and message
func New(status int, msg ...interface{}) *Error {
    message := ""
    if len(msg) == 1 {
        message = msg[0].(string)
    } else if len(msg) > 1 {
        message = fmt.Sprintf(msg[0].(string), msg[1:]...)
    }

    return &Error{status, message}
}

// Exception panic as exception
type Error struct {
    status  int
    message string
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
