package ec

import (
	"git.golaxy.org/tiny/utils/iface"
	"git.golaxy.org/tiny/utils/meta"
	"git.golaxy.org/tiny/utils/option"
	"git.golaxy.org/tiny/utils/pool"
	"git.golaxy.org/tiny/utils/uid"
)

// EntityOptions 创建实体的所有选项
type EntityOptions struct {
	CompositeFace      iface.Face[Entity]      // 扩展者，在扩展实体自身能力时使用
	Prototype          string                  // 实体原型名称
	PersistId          uid.Id                  // 实体持久化Id
	AwakeOnFirstAccess bool                    // 开启组件被首次访问时，检测并调用Awake()
	Meta               meta.Meta               // Meta信息
	ManagedPooledChunk pool.ManagedPooledChunk // 托管对象池
}

var With _Option

type _Option struct{}

// Default 默认值
func (_Option) Default() option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		With.CompositeFace(iface.Face[Entity]{})(o)
		With.Prototype("")(o)
		With.PersistId(uid.Nil)(o)
		With.AwakeOnFirstAccess(false)(o)
		With.Meta(nil)(o)
		With.ManagedPooledChunk(nil)(o)
	}
}

// CompositeFace 扩展者，在扩展实体自身能力时使用
func (_Option) CompositeFace(face iface.Face[Entity]) option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		o.CompositeFace = face
	}
}

// Prototype 实体原型名称
func (_Option) Prototype(pt string) option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		o.Prototype = pt
	}
}

// PersistId 实体持久化Id
func (_Option) PersistId(id uid.Id) option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		o.PersistId = id
	}
}

// AwakeOnFirstAccess 开启组件被首次访问时，检测并调用Awake()
func (_Option) AwakeOnFirstAccess(b bool) option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		o.AwakeOnFirstAccess = b
	}
}

// Meta Meta信息
func (_Option) Meta(m meta.Meta) option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		o.Meta = m
	}
}

// ManagedPooledChunk 托管对象池
func (_Option) ManagedPooledChunk(managed pool.ManagedPooledChunk) option.Setting[EntityOptions] {
	return func(o *EntityOptions) {
		o.ManagedPooledChunk = managed
	}
}
