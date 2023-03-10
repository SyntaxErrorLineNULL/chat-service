package repository

import "errors"

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty")

// ErrNotFound returns when chat doesn't exist
var ErrNotFound = errors.New("not found")

// ErrCannotFind returns when request was with incorrect data to search user
var ErrCannotFind = errors.New("cannot find")
