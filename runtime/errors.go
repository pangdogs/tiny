package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/utils/exception"
)

var (
	ErrContext    = fmt.Errorf("%w: runtime-context", exception.ErrTiny) // 运行时上下文错误
	ErrEntityTree = fmt.Errorf("%w: entity-tree", ErrContext)            // 实体树错误
	ErrEntityMgr  = fmt.Errorf("%w: entity-mgr", ErrContext)             // 实体管理器错误
	ErrFrame      = fmt.Errorf("%w: frame", ErrContext)                  // 帧错误
)
