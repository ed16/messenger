package domain

import "errors"

var (
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
	// ErrUserAlreadyExists will throw if the user cannot be created as i already exists
	ErrUserAlreadyExists = errors.New("user already exists")
)
