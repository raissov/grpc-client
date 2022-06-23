package models

import "errors"

var (
	ErrDBError             = errors.New("DATABASE_ERROR")
	ErrInvalidInput        = errors.New("INVALID_INPUT")
	ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
)
