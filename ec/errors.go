package ec

import (
	"fmt"
	"git.golaxy.org/tiny/utils/exception"
)

var (
	ErrEC = fmt.Errorf("%w: ec", exception.ErrTiny) // EC错误
)
