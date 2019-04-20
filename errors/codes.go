package errors

import "math"

type ErrorCode uint16

const NoCode ErrorCode = math.MaxUint16

const (
	NotFound = ErrorCode(iota)
	PermissionDenied
	PathError
	BadCredentials
	Session
	BadRequest
	ParseError
)

func (e ErrorCode) HTTPStatusCode() int {

	switch e {
	case BadRequest, ParseError:
		return 400
	case BadCredentials:
		return 401
	case PermissionDenied:
		return 403
	case NotFound, PathError:
		return 404
	default:
		return 500
	}
}
