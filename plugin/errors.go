package plugin

import (
	"fmt"
	"git.golaxy.org/tiny/utils/exception"
)

var (
	ErrPlugin = fmt.Errorf("%w: plugin", exception.ErrTiny) // 插件错误
)
