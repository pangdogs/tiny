package runtime

import (
	"git.golaxy.org/tiny/plugin"
)

// GetPluginBundle 获取插件包
func (ctx *ContextBehavior) GetPluginBundle() plugin.PluginBundle {
	return ctx.opts.PluginBundle
}
