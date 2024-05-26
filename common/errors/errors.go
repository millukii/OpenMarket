package errors

import "errors"

var (
	ErrNoItems = errors.New("items must have at least one item")
	ErrNoId = errors.New("id is required")
	ErrInvalidQuantity = errors.New("quantity is not correct")
)