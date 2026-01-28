package tiny

import (
	"fmt"

	"git.golaxy.org/tiny/utils/exception"
)

var (
	ErrTiny     = exception.ErrTiny                  // 内核错误
	ErrPanicked = exception.ErrPanicked              // panic错误
	ErrArgs     = exception.ErrArgs                  // 参数错误
	ErrRuntime  = fmt.Errorf("%w: runtime", ErrTiny) // 运行时错误
)
