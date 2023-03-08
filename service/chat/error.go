package chat

import "errors"

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty request")
