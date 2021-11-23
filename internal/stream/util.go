package stream

type HandlerFunc func(SyncItem) interface{}

func transmit(i Inlet, o Outlet) {
	defer close(i.In())
	for ele := range o.Out() {
		i.In() <- ele
	}
}

/*func DoStream(i Inlet, o Outlet) {
	go transmit(i, o)
}*/
