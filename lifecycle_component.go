package tiny

type LifecycleComponentAwake interface {
	Awake()
}

type LifecycleComponentStart interface {
	Inited()
}

type LifecycleComponentUpdate = eventUpdate

type LifecycleComponentLateUpdate = eventLateUpdate

type LifecycleComponentShut interface {
	Shut()
}
