package repository

import "errors"

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty")

// ErrNotFound returns when record doesn't exist in database
var ErrNotFound = errors.New("not found")

// ErrCannotFind returns when request was with incorrect data to search user
var ErrCannotFind = errors.New("cannot find")
