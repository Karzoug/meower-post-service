package ucerr

import (
	"github.com/Karzoug/meower-post-service/pkg/ucerr/codes"
)

// Error is a service/usecase level error.
type Error struct {
	msg  string
	err  error
	code codes.Code
}

func NewError(err error, msg string, code codes.Code) Error {
	return Error{
		msg:  msg,
		err:  err,
		code: code,
	}
}

func NewInternalError(err error) Error {
	return Error{
		msg:  "Internal error",
		err:  err,
		code: codes.Internal,
	}
}

// Error returns error message which can be returned to the client.
func (e Error) Error() string {
	return e.msg
}

func (e Error) Code() codes.Code {
	return e.code
}

func (e Error) Unwrap() error {
	return e.err
}
