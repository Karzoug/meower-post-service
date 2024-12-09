package repo

import "errors"

var (
	ErrNoAffected          = errors.New("no affected")
	ErrAborted             = errors.New("aborted")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrRecordNotFound      = errors.New("record not found")
)
