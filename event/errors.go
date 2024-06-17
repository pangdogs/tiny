package event

import (
	"fmt"
	"git.golaxy.org/tiny/utils/exception"
)

var (
	ErrEvent = fmt.Errorf("%w: event", exception.ErrTiny) // 事件错误
	ErrArgs  = exception.ErrArgs                          // 参数错误
)
