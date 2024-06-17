package runtime

import (
	"git.golaxy.org/tiny/utils/async"
	"git.golaxy.org/tiny/utils/uid"
)

// Deprecated: UnsafeContext 访问运行时上下文内部方法
func UnsafeContext(ctx Context) _UnsafeContext {
	return _UnsafeContext{
		Context: ctx,
	}
}

type _UnsafeContext struct {
	Context
}

// Init 初始化
func (uc _UnsafeContext) Init(opts ContextOptions) {
	uc.Context.init(opts)
}

// GetOptions 获取运行时上下文所有选项
func (uc _UnsafeContext) GetOptions() *ContextOptions {
	return uc.getOptions()
}

// NewId 创建Id
func (uc _UnsafeContext) NewId() uid.Id {
	return uc.newId()
}

// SetFrame 设置帧
func (uc _UnsafeContext) SetFrame(frame Frame) {
	uc.setFrame(frame)
}

// SetCallee 设置调用接受者
func (uc _UnsafeContext) SetCallee(callee async.Callee) {
	uc.setCallee(callee)
}

// ChangeRunningState 修改运行状态
func (uc _UnsafeContext) ChangeRunningState(state RunningState) {
	uc.changeRunningState(state)
}

// GC GC
func (uc _UnsafeContext) GC() {
	uc.gc()
}
