package devices

import "errors"

var (
	ErrInvalidAction        = errors.New("invalid action requested")
	ErrSceneSystemError     = errors.New("scene system error")
	ErrBadStateValueType    = errors.New("bad state value type")
	ErrNotSupportedAbsValue = errors.New("absolute value not supported for this device")
)
