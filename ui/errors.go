package ui

import "errors"

var (
	errNoEntitySelected   = errors.New("no entity selected")
	errYouDontExist       = errors.New("you don't exist")
	errAlreadyHereWarning = errors.New("you are already here")
)
