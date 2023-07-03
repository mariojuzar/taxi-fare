package fare

import "errors"

var (
	ErrImproperFormat      = errors.New("improper input format")
	ErrBlankLine           = errors.New("receive a blank input")
	ErrPastTime            = errors.New("receive past time")
	ErrIntervalTimeTooLong = errors.New("time interval is more than 5 minutes from latest data")
	ErrNotEnoughData       = errors.New("records is less than 2")
	ErrZeroMileage         = errors.New("invalid zero movement mileage")
	ErrInvalidStart        = errors.New("invalid start position")
)
