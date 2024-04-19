package storage

import "errors"

var (
	ErrorUserAlreadyExist = errors.New("user already exist")
	ErrorUserNotFound     = errors.New("user not found")
	ErrorAppNotFound      = errors.New("app not found")
)
