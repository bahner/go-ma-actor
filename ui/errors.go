package ui

import "errors"

var (
	ErrNoEntitySelected    = errors.New("no entity selected")
	ErrYouDontExist        = errors.New("you don't exist")
	ErrFailedToEnterEntity = errors.New("failed to enter entity")
	ErrSelfEntryWarning    = errors.New("entering yourself")
	ErrAlreadyHereWarning  = errors.New("you are already here")
)
