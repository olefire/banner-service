package repository

import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrAlreadyExists  = errors.New("record already exists")
	ErrBannerInactive = errors.New("banner is inactive")
)
