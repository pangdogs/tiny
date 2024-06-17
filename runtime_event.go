//go:generate go run git.golaxy.org/tiny/event/eventcode gen_event --default_export=false --default_auto=false
package tiny

type eventUpdate interface {
	Update()
}

type eventLateUpdate interface {
	LateUpdate()
}
