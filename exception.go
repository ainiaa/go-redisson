package redis

import (
	"github.com/ainiaa/go-exception"
)

var (
	ErrConnTypeUnknown = exception.New(-1, "unknown connect type")
	ErrPing            = exception.New(-2, "redis conn error")
	ErrNoReady         = exception.New(-3, "redis is not ready,please init it first")
)
var (
	ErrExitsLock           = exception.New(-4, "lock exits")
	ErrAcquiredLock        = exception.New(-5, "acquire lock error")
	ErrAcquiredLockTimeout = exception.New(-6, "acquire lock timeout error")
)
