package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested data is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict     = errors.New("Your data already exist")
	ErrAccountExist = errors.New("Account already exist")
	ErrEmailExist   = errors.New("Email already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput   = errors.New("Given Param is not valid")
	ErrPassword        = errors.New("Wrong Password")
	ErrEmailNotFound   = errors.New("Email Not Found")
	ErrorAuthorization = errors.New("Unathorized")
)
