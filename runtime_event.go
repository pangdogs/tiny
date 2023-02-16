//go:generate go run kit.golaxy.org/tiny/localevent/eventcode --decl_file=$GOFILE gen_emit --package=$GOPACKAGE --default_export=0
package tiny

type eventUpdate interface {
	Update()
}

type eventLateUpdate interface {
	LateUpdate()
}
