package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested data is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your data already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")

	// account
	ErrAccountExist       = errors.New("account already exist")
	ErrEmailExist         = errors.New("email already exist")
	ErrPassword           = errors.New("wrong Password")
	ErrEmailNotFound      = errors.New("email Not Found")
	ErrorAuthorization    = errors.New("unathorized")
	ErrorEmailNotVerified = errors.New("email not verified")

	// generateUrl
	ErrUrlNotFound       = errors.New("url not found")
	ErrUrlOriginExist    = errors.New("url origin already exist")
	ErrUrlGeneratedExist = errors.New("url generated already exist")
	ErrNameIsExist       = errors.New("name is exist")
)
