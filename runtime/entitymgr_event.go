//go:generate go run kit.golaxy.org/tiny/localevent/eventcode --decl_file=$GOFILE gen_emit --package=$GOPACKAGE
package runtime

import "kit.golaxy.org/tiny/ec"

// EventEntityMgrAddEntity [EmitUnExport] 事件：实体管理器添加实体
type EventEntityMgrAddEntity interface {
	OnEntityMgrAddEntity(entityMgr IEntityMgr, entity ec.Entity)
}

// EventEntityMgrRemovingEntity [EmitUnExport] 事件：实体管理器开始删除实体
type EventEntityMgrRemovingEntity interface {
	OnEntityMgrRemovingEntity(entityMgr IEntityMgr, entity ec.Entity)
}

// EventEntityMgrRemoveEntity [EmitUnExport] 事件：实体管理器删除实体
type EventEntityMgrRemoveEntity interface {
	OnEntityMgrRemoveEntity(entityMgr IEntityMgr, entity ec.Entity)
}

// EventEntityMgrEntityAddComponents [EmitUnExport] 事件：实体管理器中的实体添加组件
type EventEntityMgrEntityAddComponents interface {
	OnEntityMgrEntityAddComponents(entityMgr IEntityMgr, entity ec.Entity, components []ec.Component)
}

// EventEntityMgrEntityRemoveComponent [EmitUnExport] 事件：实体管理器中的实体删除组件
type EventEntityMgrEntityRemoveComponent interface {
	OnEntityMgrEntityRemoveComponent(entityMgr IEntityMgr, entity ec.Entity, component ec.Component)
}
