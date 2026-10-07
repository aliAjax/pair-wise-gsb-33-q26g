package repository

import "errors"

// Sentinel errors shared by all repositories.
var (
	ErrNotFound   = errors.New("record not found")
	ErrDuplicate  = errors.New("duplicate record")
	ErrConflict   = errors.New("conflict state")
)
