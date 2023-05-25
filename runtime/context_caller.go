package runtime

import (
	"fmt"
	"kit.golaxy.org/tiny/ec"
	"kit.golaxy.org/tiny/internal"
)

// NewRet 创建调用结果
var NewRet = internal.NewRet

// Ret 调用结果
type Ret = internal.Ret

// AsyncRet 异步调用结果
type AsyncRet = internal.AsyncRet

// Caller 异步调用发起者
type Caller = internal.Caller

// Callee 调用接收者
type Callee interface {
	// PushCall 将代码片段压入接收者的任务处理流水线，串行化的进行调用。
	PushCall(segment func())
}

func entityCaller(entity ec.Entity) Caller {
	return Get(entity)
}

func entityExist(entity ec.Entity) bool {
	_, ok := Get(entity).GetEntityMgr().GetEntity(entity.GetId())
	return ok
}

// SyncCall 同步调用。在运行时中，将代码片段压入任务流水线，串行化的进行调用，会阻塞并等待返回值。
//
//	注意：
//	- 代码片段中的线程安全问题。
//	- 当运行时的SyncCallTimeout选项设置为0时，在代码片段中，如果向调用方所在的运行时发起同步调用，那么会造成线程死锁。
//	- 调用过程中的panic信息，均会转换为error返回。
func (ctx *ContextBehavior) SyncCall(segment func() Ret) Ret {
	var ret Ret

	func() {
		defer func() {
			if info := recover(); info != nil {
				err, ok := info.(error)
				if !ok {
					err = fmt.Errorf("%v", info)
				}
				ret = NewRet(err, nil)
			}
		}()

		if segment == nil {
			panic("nil segment")
		}

		ctx.callee.PushCall(func() {
			ret = segment()
		})
	}()

	return ret
}

// AsyncCall 异步调用。在运行时中，将代码片段压入任务流水线，串行化的进行调用，不会阻塞，会返回AsyncRet。
//
//	注意：
//	- 代码片段中的线程安全问题。
//	- 在代码片段中，如果向调用方所在的运行时发起同步调用，并且调用方也在阻塞AsyncRet等待返回值，那么会造成线程死锁。
//	- 调用过程中的panic信息，均会转换为error返回。
func (ctx *ContextBehavior) AsyncCall(segment func() Ret) AsyncRet {
	asyncRet := make(chan Ret, 1)

	go func() {
		defer func() {
			if info := recover(); info != nil {
				err, ok := info.(error)
				if !ok {
					err = fmt.Errorf("%v", info)
				}
				asyncRet <- NewRet(err, nil)
				close(asyncRet)
			}
		}()

		if segment == nil {
			panic("nil segment")
		}

		ctx.callee.PushCall(func() {
			defer close(asyncRet)
			asyncRet <- segment()
		})
	}()

	return asyncRet
}

// SyncCallNoRet 同步调用，无返回值。在运行时中，将代码片段压入任务流水线，串行化的进行调用，会阻塞，没有返回值。
//
//	注意：
//	- 代码片段中的线程安全问题。
//	- 当运行时的SyncCallTimeout选项设置为0时，在代码片段中，如果向调用方所在的运行时发起同步调用，那么会造成线程死锁。
//	- 调用过程中的panic信息，均会抛弃。
func (ctx *ContextBehavior) SyncCallNoRet(segment func()) {
	if segment == nil {
		return
	}

	func() {
		defer func() {
			recover()
		}()

		ctx.callee.PushCall(segment)
	}()
}

// AsyncCallNoRet 异步调用，无返回值。在运行时中，将代码片段压入任务流水线，串行化的进行调用，不会阻塞，没有返回值。
//
//	注意：
//	- 代码片段中的线程安全问题。
//	- 调用过程中的panic信息，均会抛弃。
func (ctx *ContextBehavior) AsyncCallNoRet(segment func()) {
	if segment == nil {
		return
	}

	go func() {
		defer func() {
			recover()
		}()

		ctx.callee.PushCall(segment)
	}()
}
