package cacheerr

import "errors"

var (
	ErrMissingHost = errors.New("missing redis host")
	ErrInvalidPort = errors.New("redis port must be between 1 and 65535")
	ErrInvalidDB   = errors.New("redis DB index must be >= 0")

	ErrKeyNotFound    = errors.New("key not found in cache")
	ErrDecode         = errors.New("failed to decode to struct")
	ErrSetNil         = errors.New("cannot cache nil value")
	ErrSerializeValue = errors.New("failed to serialize cache value")
)
