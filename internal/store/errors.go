package store

import "errors"

// ErrNotFound is returned when the requested env set does not exist on disk.
var ErrNotFound = errors.New("env set not found")

// ErrBadPassphrase is returned when decryption fails, likely due to a wrong passphrase.
var ErrBadPassphrase = errors.New("incorrect passphrase or corrupted data")
