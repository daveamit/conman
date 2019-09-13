package conman

import "errors"

// ErrNotImplemented :- Method / Function / Functionality no implemented
var ErrNotImplemented = errors.New("Method / Function / Functionality no implemented")

// ErrAlreadyWatchingGivenKey :- Method / Function / Functionality no implemented
var ErrAlreadyWatchingGivenKey = errors.New("Already watching given key")

// ErrKeyNotFound :- the key does not exist
var ErrKeyNotFound = errors.New("Key you are trying to access does not exist")
