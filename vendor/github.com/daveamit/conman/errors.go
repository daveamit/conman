package conman

import "errors"

// ErrNotImplemented :- Method / Function / Functionality no implemented
var ErrNotImplemented = errors.New("Method / Function / Functionality no implemented")

// ErrAlreadyWatchingGivenKey :- Method / Function / Functionality no implemented
var ErrAlreadyWatchingGivenKey = errors.New("Already watching given key")
